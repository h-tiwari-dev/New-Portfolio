package helpers

import (
	"app/db"
	"app/models"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	md_html "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

func init() {
	htmlFormatter = html.New(html.WithClasses(false), html.TabWidth(4))
	if htmlFormatter == nil {
		panic("couldn't create html formatter")
	}
	styleName := "monokailight"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		panic(fmt.Sprintf("didn't find style '%s'", styleName))
	}
}

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return htmlFormatter.Format(w, highlightStyle, it)
}

// an actual rendering of Paragraph is more complicated
func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)

	// TODO: Remove this
	//err := htmlFormatter.WriteCSS(w, highlightStyle)
	//if err != nil {
	//	panic("Formatting gone wrong")
	//}
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func modifyAstToSaveImageSrcs(doc ast.Node, strict bool) ([]uint, ast.Node) {
	var image_ids []uint

	// Walk the AST
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if img, ok := node.(*ast.Image); ok && entering {
			// Extract the image source URL
			src := string(img.Destination)

			imageContent, err := os.ReadFile(fmt.Sprintf("./blogs/md/%s", strings.TrimPrefix(src, "./")))
			if err != nil && strict {
				panic(fmt.Sprintf("StrictMODE:- Failed to read image file: %v \nSwitch off the strict mode for more relax conversion.", err))
			}

			imageData := models.Image{
				Filename: "deafult_img",
				Data:     []byte(imageContent),
				TopImage: 0,
			}

			attr := img.Attribute
			if attr == nil {
				attr = &ast.Attribute{}
				attr.Attrs = make(map[string][]byte)
			}

			if err := db.DB.Create(&imageData).Error; err != nil {
				if strict {
					panic(
						fmt.Sprintf(
							"StrictMODE:- Failed to save image to db: %v %s \nSwitch off the strict mode for more relax conversion.",
							err,
							imageData.Filename,
						),
					)
				}
			} else {
				attr.Attrs["hx-get"] = []byte(fmt.Sprintf("/image/%d", imageData.ID))
				attr.Attrs["hx-trigger"] = []byte(`load`)
				attr.Attrs["hx-target"] = []byte(`this`)
				attr.Attrs["hx-swap"] = []byte(`outerHTML`)
				img.Destination = []byte(``)

				image_ids = append(image_ids, imageData.ID)
			}

			img.Attribute = attr
		}

		// Continue walking the AST
		return ast.GoToNext
	})

	// Return the modified AST and the list of image source URLs
	return image_ids, doc
}

func newCustomizedRenderHook() *md_html.Renderer {
	opts := md_html.RendererOptions{
		Flags:          md_html.CommonFlags | md_html.HrefTargetBlank,
		RenderNodeHook: myRenderHook,
	}
	return md_html.NewRenderer(opts)
}

func ConvertMdToHTML(content []byte) ([]uint, []byte) {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	imageSrcs, modifiedDoc := modifyAstToSaveImageSrcs(doc, true)

	renderer := newCustomizedRenderHook()
	html := markdown.Render(modifiedDoc, renderer)

	return imageSrcs, html
}
