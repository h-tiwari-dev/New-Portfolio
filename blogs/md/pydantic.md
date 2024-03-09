# Elevate Your LLM's Game with Pydantic! Part I

In the world of **Large Language Models (LLMs)**, making them play nice with our applications is key. We want our models to dish out neat JSON for easy integration into our full-stack setups. But relying solely on LLM luck for perfect JSON? Let's be realistic.

Suppose you're incorporating an LLM into your app, striving for precise JSON output. Considering the importance of this data, we might need to save it for the next steps in our logic. You provide clear prompts, cross your fingers, and hope. Yet, hope isn't a strategy, and guarantees are scarce.

Meet [**Pydantic**](https://pypi.org/project/pydantic/), a handy data validation tool. This tool turn your JSON into a structured class for order in the chaos. Plus, Pydantic brings validations and extra functionality to the table.

We'll also use [**Instructor**](https://pypi.org/project/instructor/). **Instructor** patches our OpenAI client, empowering it to return our response model (essentially our Pydantic class).Additionally, we can incorporate **Max Retries** to automatically retry when our LLM fails to deliver the desired output


In simple terms, Pydantic takes your JSON, turns it into a class, and lets you add checks and tweaks. The real win? When Pydantic teams up with LLMs, making your applications more reliable and functional. Get ready to supercharge your LLM game with Pydantic magic!


# Tackling the JSON Conundrum

Reliability in our outputs is really important (else we'll see you in try-catch hell). Consider this scenario: imagine you're crafting a medical application tasked with extracting information from a string. The next logical step is to convert this data into JSON for further analysis. Perhaps, you plan to map this JSON to a class object, storing it either temporarily in memory or persisting it in a database, say using SQLAlchemy.

In this process, the challenge lies in ensuring that the JSON output remains accurate and consistent, ready to be seamlessly integrated into your application's logic. This is where the crux of the problem resides.

Let's take our medical example and flesh it out.

Suppose we want this information :- 
```python
medical_info = """Sex: Female, Age: 79
Geographical region: North America
Pathology: Spontaneous pneumothorax
Symptoms:
---------
 - I have chest pain even at rest.
 - I feel pain.
 - The pain is:
     » a knife stroke
 - The pain locations are:
     » upper chest
     » breast(R)
     » breast(L)
 - On a scale of 0-10, the pain intensity is 7
 - On a scale of 0-10, the pain's location precision is 4
 - On a scale of 0-10, the pace at which the pain appear is 9
 - I have symptoms that increase with physical exertion but alleviate with rest.
Antecedents:
-----------
 - I have had a spontaneous pneumothorax.
 - I smoke cigarettes.
 - I have a chronic obstructive pulmonary disease.
 - Some family members have had a pneumothorax.
Differential diagnosis:
----------------------
```

# Elevate Your LLM's Game with Pydantic! Part I

In the world of **Large Language Models (LLMs)**, making them play nice with our applications is key. We want our models to dish out neat JSON for easy integration into our full-stack setups. But relying solely on LLM luck for perfect JSON? Let's be realistic.

Suppose you're incorporating an LLM into your app, striving for precise JSON output. Considering the importance of this data, we might need to save it for the next steps in our logic. You provide clear prompts, cross your fingers, and hope. Yet, hope isn't a strategy, and guarantees are scarce.

Meet [**Pydantic**](https://pypi.org/project/pydantic/), a handy data validation tool. This tool turn your JSON into a structured class for order in the chaos. Plus, Pydantic brings validations and extra functionality to the table.

We'll also use [**Instructor**](https://pypi.org/project/instructor/). **Instructor** patches our OpenAI client, empowering it to return our response model (essentially our Pydantic class).Additionally, we can incorporate **Max Retries** to automatically retry when our LLM fails to deliver the desired output


In simple terms, Pydantic takes your JSON, turns it into a class, and lets you add checks and tweaks. The real win? When Pydantic teams up with LLMs, making your applications more reliable and functional. Get ready to supercharge your LLM game with Pydantic magic!


# Tackling the JSON Conundrum

Reliability in our outputs is really important (else we'll see you in try-catch hell). Consider this scenario: imagine you're crafting a medical application tasked with extracting information from a string. The next logical step is to convert this data into JSON for further analysis. Perhaps, you plan to map this JSON to a class object, storing it either temporarily in memory or persisting it in a database, say using SQLAlchemy.

In this process, the challenge lies in ensuring that the JSON output remains accurate and consistent, ready to be seamlessly integrated into your application's logic. This is where the crux of the problem resides.

Let's take our medical example and flesh it out.

Suppose we want this information :- 
```python
medical_info = """Sex: Female, Age: 79
    Geographical region: North America
    Pathology: Spontaneous pneumothorax
    Symptoms:
    ---------
    - I have chest pain even at rest.
    - I feel pain.
    - The pain is:
    » a knife stroke
    - The pain locations are:
        » upper chest
        » breast(R)
        » breast(L)
        - On a scale of 0-10, the pain intensity is 7
        - On a scale of 0-10, the pain's location precision is 4
        - On a scale of 0-10, the pace at which the pain appear is 9
        - I have symptoms that increase with physical exertion but alleviate with rest.
        Antecedents:
        -----------
        - I have had a spontaneous pneumothorax.
        - I smoke cigarettes.
        - I have a chronic obstructive pulmonary disease.
        - Some family members have had a pneumothorax.
        Differential diagnosis:
        ----------------------
        Unstable angina: 0.262, Stable angina: 0.201, Possible NSTEMI / STEMI: 0.160, GERD: 0.145, Pericarditis: 0.091, Atrial fibrillation: 0.082, Spontaneous pneumothorax: 0.060
"""
```

To be converted into this format:-

```python
    json_format = """{
        "patient_info": {
            "sex": "",
                "age": ,
                "geographical_region": "",
        },
            "medical_history": {
                "pathology": "",
                "symptoms": {
                    "description": "",
                    "pain": {
                        "type": "",
                        "locations": [],
                        "intensity": ,
                        "location_precision": ,
                        "pace": ,
                    },
                    "increase_with_exertion": true/false,
                    "alleviate_with_rest": true/false,
                },
            },
            "risk_factors": {},
            "differential_diagnosis": [
            {
                "disease_name": "",
                "probability": 
            },
            ]
    }"""

```
> Let's take a sec and think what we are doing. If we assume our llms as a black box:-
> ```python
>    def llm(prompt: str, schema: str) -> str:
>       pass  # Black Magic, and hope to receive valid json.
>    ```

***Now let's plead to the AI goddess to convert this into valid JSON.***

```python
completion = openai_client.chat.completions.create(
        model="gpt-3.5-turbo",
        messages=[{
        "role": "user", 
        "content": f"Please convert the following information into valid json representing the medical diagnosis ${medical_info}. Please convert the data in the following format and fill in the data ${json_format}"
        }]
        )
# Now let's extract our "valid" json
dump = completion.model_dump()
    medical_info  =json.loads(dump["choices"][0]["message"]["content"])
    print(json.dumps(medical_info, indent=2)) # A big leap of faith.
```


In the code, relying solely on json.dumps might lead to errors if the model doesn't provide valid JSON. Adding a try-except block for error handling and incorporating a retry mechanism can be quite cumbersome. Dealing with these uncertainties emphasizes the challenges of ensuring a smooth interaction with the language model output.
By the way output we get looks something like this

```json
{
    "patient_info": {
        "sex": "Female",
        "age": 79,
        "geographical_region": "North America"
    },
        "medical_history": {
            "pathology": "Spontaneous pneumothorax",
            "symptoms": {
                "description": "I have chest pain even at rest. I feel pain.",
                "pain": {
                    "type": "a knife stroke",
                    "locations": [
                        "upper chest",
                    "breast(R)",
                    "breast(L)"
                    ],
                    "intensity": 7,
                    "location_precision": 4,
                    "pace": 9
                },
                "increase_with_exertion": true,
                "alleviate_with_rest": true
            }
        },
        ...
            "probability": 0.06
}
]
}
```

The issue of uncertainties becomes more pronounced when dealing with complex data structures or interconnected structures.


We're crossing our fingers, hoping that when we convert our LLM output, a string supposedly in valid JSON format, into our object, everything works smoothly. However, in our current testing example, a couple of issues are still lingering:

**1. Lack of Type Safety:**
The current approach involves converting a string to a JSON object, and we're essentially relying on the all-powerful AI god to provide us with correct JSON. What if, instead of a birthdate, we need...

**2. Validation Issues:**
Handling input validation manually is a bit of a headache. To validate, we have to manually check the structure of the JSON, which results in a messy function like this:

```python
def validate_json_structure(json_string):
    try:
    data = json.loads(json_string)

# Validate patient_info
    patient_info = data.get("patient_info")
    if not patient_info or not isinstance(patient_info, dict):
        return False

# Validate sex, age, and geographical_region in patient_info
#... os on with more validations within validations
        return True

        except json.JSONDecodeError:
        return False
        ```

        What a horrible mess. It's not the most elegant solution. (Psst, we'll soon explore how Pydantic can simplify this mess and add various validations.)

        On another note, Pydantic allows us to chain our prompts using inheritance(OOP), as you'll see in an example towards the end of this blog.

## A Structured Approach with Pydantic and Instructor.
        We aim for our magical function to receive a schema defined as a Python class or model and return either the same or another class/model. It should look something like this:-

        ```python
        def llm(prompt: str, schema: Model) -> Model:
        pass  
        ```

        This is where Pydantic steps in. Let's import the necessary modules and set up our OpenAI client with the help of [Instructor](https://pypi.org/project/instructor/):

        ```python
        import instructor

        instructor_openai_client = instructor.patch(openai.Client(
                    api_key=open_ai_key, organization=open_ai_org_key, timeout=20000, max_retries=3
                    ))
        ```

        Overall, Instructor is a user-friendly, transparent, and Pythonic solution for leveraging OpenAI's function calling to extract data. It patches to the OpenAI's library and help us achive the ```(prompt, model) -> model``` structure.

        ---

        Next, we define our JSON structure using Pydantic classes. This approach allows us to include additional docstrings for field descriptions and other useful information. All of this aids the language model in generating or extracting information from the context provided by the model.

        ```python
        from typing import List, Literal, Optional
        from pydantic import BaseModel, Field

        class Symptoms(BaseModel):
            """
                Represents the symptoms of a patient.
                """
                description: str = Field(description="A general scientific and objective description of the symptoms.")
                pain_type: str
                locations: List[str]
                intensity: int
                location_precision: int
                pace: int

                class MedicalHistory(BaseModel):
                    pathology: str
                    symptoms: Symptoms
                    increase_with_exertion: bool
                    alleviate_with_rest: bool

                    class RiskFactors(BaseModel):
                        spontaneous_history: bool
                        smoking_history: bool
                        copd_history: bool
                        family_history: str

                        class DifferentialDiagnosis(BaseModel):
                            disease_name: str
                            probability: float

                            class PatientInfo(BaseModel):
                                sex: Literal['M', 'F']
                                age: int
                                geographical_region: str

                                class PatientData(BaseModel):
                                    patient_info: PatientInfo
                                    medical_history: MedicalHistory
                                    risk_factors: RiskFactors
                                    differential_diagnosis: List[DifferentialDiagnosis]
                                    ```

### Now, let's utilize our magical function with the OpenAI completion
                                    ```python
                                    completion = instructor_openai_client.chat.completions.create(
                                            model="gpt-3.5-turbo",
                                            messages=[{"role": "user", "content": f"Please convert the following information into valid JSON representing the medical diagnosis {medical_info}."}],
                                            response_model=MedicalInfo  # Replace with the appropriate response model
                                            )
    print(type(completion))
print(json.dumps(completion.model_dump(), indent=1))
    ```
    **Voila**:-
    ```shell
    <class '__main__.PatientData'> # Notice how the type of the data structure we got is a class!!!
{
    "patient_info": {
        "sex": "F",
        "age": 79,
        "geographical_region": "North America"
    },
        "medical_history": {
            "pathology": "Spontaneous pneumothorax",
            "symptoms": {
                "description": "I have chest pain even at rest. I feel pain. The pain is a knife stroke. The pain locations are upper chest, breast(R), breast(L). On a scale of 0-10, the pain intensity is 7. On a scale of 0-10, the pain's location precision is 4. On a scale of 0-10, the pace at which the pain appears is 9. I have symptoms that increase with physical exertion but alleviate with rest.",
                "pain_type": "knife stroke",
                "locations": [
                    "upper chest",
                "breast(R)",
                "breast(L)"
                ],
                "intensity": 7,
                "location_precision": 4,
                "pace": 9
            },
            "increase_with_exertion": true,
            "alleviate_with_rest": true
        },
        "risk_factors": {
            ...
                "probability": 0.06
        }
    ]
}
```

By setting `response_model` to `MedicalInfo`, we ensure a clear output structure. Pydantic guarantees data adherence, streamlining integration and providing a type hint of `PatientData`.

Pydantic organizes JSON with automatic validation. Deviations trigger validation errors, ensuring data integrity.

Docstrings and field descriptions aid developers and shape the JSON schema for OpenAI. Navigate confidently with structured, validated data, and notice the response type as `PatientData` for seamless integration.

# **Congratulations, It's a class!**

In the next part of this series, We'll talk about LLM validations, seamless retry mechanisms, how you can create complex data structures like directed acyclic graphs (DAGs), and much more using Pydantic. Stay tuned for the next part.

> References:-
> 1. This blog post is inspired by an awesome talk by [Jason Liu](https://github.com/jxnl)'s, [talk](https://www.youtube.com/watch?v=yj-wSRJwrrc). Please watch it too for better reference.
> 2. [Pydantic](https://pypi.org/project/pydantic/), [Instructor](https://pypi.org/project/instructor/)


nstable angina: 0.262, Stable angina: 0.201, Possible NSTEMI / STEMI: 0.160, GERD: 0.145, Pericarditis: 0.091, Atrial fibrillation: 0.082, Spontaneous pneumothorax: 0.060
"""
```
To be converted into this format:-

```python
json_format = """{
    "patient_info": {
        "sex": "",
        "age": ,
        "geographical_region": "",
    },
    "medical_history": {
        "pathology": "",
        "symptoms": {
            "description": "",
            "pain": {
                "type": "",
                "locations": [],
                "intensity": ,
                "location_precision": ,
                "pace": ,
            },
            "increase_with_exertion": true/false,
            "alleviate_with_rest": true/false,
        },
    },
    "risk_factors": {},
    "differential_diagnosis": [
    {
        "disease_name": "",
        "probability": 
    },
]
}"""
```
> Let's take a sec and think what we are doing. If we assume our llms as a black box:-
> ```python
>    def llm(prompt: str, schema: str) -> str:
>       pass  # Black Magic, and hope to receive valid json.
>    ```

***Now let's plead to the AI goddess to convert this into valid JSON.***

```python
completion = openai_client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{
        "role": "user", 
        "content": f"Please convert the following information into valid json representing the medical diagnosis ${medical_info}. Please convert the data in the following format and fill in the data ${json_format}"
    }]
)
# Now let's extract our "valid" json
dump = completion.model_dump()
medical_info  =json.loads(dump["choices"][0]["message"]["content"])
print(json.dumps(medical_info, indent=2)) # A big leap of faith.
```

In the code, relying solely on json.dumps might lead to errors if the model doesn't provide valid JSON. Adding a try-except block for error handling and incorporating a retry mechanism can be quite cumbersome. Dealing with these uncertainties emphasizes the challenges of ensuring a smooth interaction with the language model output.

By the way output we get looks something like this
```json
{
  "patient_info": {
    "sex": "Female",
    "age": 79,
    "geographical_region": "North America"
  },
  "medical_history": {
    "pathology": "Spontaneous pneumothorax",
    "symptoms": {
      "description": "I have chest pain even at rest. I feel pain.",
      "pain": {
        "type": "a knife stroke",
        "locations": [
          "upper chest",
          "breast(R)",
          "breast(L)"
        ],
        "intensity": 7,
        "location_precision": 4,
        "pace": 9
      },
      "increase_with_exertion": true,
      "alleviate_with_rest": true
    }
  },
...
      "probability": 0.06
    }
  ]
}
```

The issue of uncertainties becomes more pronounced when dealing with complex data structures or interconnected structures.


We're crossing our fingers, hoping that when we convert our LLM output, a string supposedly in valid JSON format, into our object, everything works smoothly. However, in our current testing example, a couple of issues are still lingering:

**1. Lack of Type Safety:**
   The current approach involves converting a string to a JSON object, and we're essentially relying on the all-powerful AI god to provide us with correct JSON. What if, instead of a birthdate, we need...

**2. Validation Issues:**
   Handling input validation manually is a bit of a headache. To validate, we have to manually check the structure of the JSON, which results in a messy function like this:

```python
def validate_json_structure(json_string):
    try:
        data = json.loads(json_string)

        # Validate patient_info
        patient_info = data.get("patient_info")
        if not patient_info or not isinstance(patient_info, dict):
            return False

        # Validate sex, age, and geographical_region in patient_info
        #... os on with more validations within validations
        return True

    except json.JSONDecodeError:
        return False
```

What a horrible mess. It's not the most elegant solution. (Psst, we'll soon explore how Pydantic can simplify this mess and add various validations.)

On another note, Pydantic allows us to chain our prompts using inheritance(OOP), as you'll see in an example towards the end of this blog.

## A Structured Approach with Pydantic and Instructor.
We aim for our magical function to receive a schema defined as a Python class or model and return either the same or another class/model. It should look something like this:-

```python
def llm(prompt: str, schema: Model) -> Model:
    pass  
```

This is where Pydantic steps in. Let's import the necessary modules and set up our OpenAI client with the help of [Instructor](https://pypi.org/project/instructor/):

```python
import instructor

instructor_openai_client = instructor.patch(openai.Client(
    api_key=open_ai_key, organization=open_ai_org_key, timeout=20000, max_retries=3
))
```

Overall, Instructor is a user-friendly, transparent, and Pythonic solution for leveraging OpenAI's function calling to extract data. It patches to the OpenAI's library and help us achive the ```(prompt, model) -> model``` structure.

---

Next, we define our JSON structure using Pydantic classes. This approach allows us to include additional docstrings for field descriptions and other useful information. All of this aids the language model in generating or extracting information from the context provided by the model.

```python
from typing import List, Literal, Optional
from pydantic import BaseModel, Field

class Symptoms(BaseModel):
    """
        Represents the symptoms of a patient.
    """
    description: str = Field(description="A general scientific and objective description of the symptoms.")
    pain_type: str
    locations: List[str]
    intensity: int
    location_precision: int
    pace: int

class MedicalHistory(BaseModel):
    pathology: str
    symptoms: Symptoms
    increase_with_exertion: bool
    alleviate_with_rest: bool

class RiskFactors(BaseModel):
    spontaneous_history: bool
    smoking_history: bool
    copd_history: bool
    family_history: str

class DifferentialDiagnosis(BaseModel):
    disease_name: str
    probability: float

class PatientInfo(BaseModel):
    sex: Literal['M', 'F']
    age: int
    geographical_region: str
            
class PatientData(BaseModel):
    patient_info: PatientInfo
    medical_history: MedicalHistory
    risk_factors: RiskFactors
    differential_diagnosis: List[DifferentialDiagnosis]
```

### Now, let's utilize our magical function with the OpenAI completion
```python
completion = instructor_openai_client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": f"Please convert the following information into valid JSON representing the medical diagnosis {medical_info}."}],
    response_model=MedicalInfo  # Replace with the appropriate response model
)
print(type(completion))
print(json.dumps(completion.model_dump(), indent=1))
```
**Voila**:-
```shell
<class '__main__.PatientData'> # Notice how the type of the data structure we got is a class!!!
{
 "patient_info": {
  "sex": "F",
  "age": 79,
  "geographical_region": "North America"
 },
 "medical_history": {
  "pathology": "Spontaneous pneumothorax",
  "symptoms": {
   "description": "I have chest pain even at rest. I feel pain. The pain is a knife stroke. The pain locations are upper chest, breast(R), breast(L). On a scale of 0-10, the pain intensity is 7. On a scale of 0-10, the pain's location precision is 4. On a scale of 0-10, the pace at which the pain appears is 9. I have symptoms that increase with physical exertion but alleviate with rest.",
   "pain_type": "knife stroke",
   "locations": [
    "upper chest",
    "breast(R)",
    "breast(L)"
   ],
   "intensity": 7,
   "location_precision": 4,
   "pace": 9
  },
  "increase_with_exertion": true,
  "alleviate_with_rest": true
 },
 "risk_factors": {
...
   "probability": 0.06
  }
 ]
}
```

By setting `response_model` to `MedicalInfo`, we ensure a clear output structure. Pydantic guarantees data adherence, streamlining integration and providing a type hint of `PatientData`.

Pydantic organizes JSON with automatic validation. Deviations trigger validation errors, ensuring data integrity.

Docstrings and field descriptions aid developers and shape the JSON schema for OpenAI. Navigate confidently with structured, validated data, and notice the response type as `PatientData` for seamless integration.

# **Congratulations, It's a class!**

In the next part of this series, We'll talk about LLM validations, seamless retry mechanisms, how you can create complex data structures like directed acyclic graphs (DAGs), and much more using Pydantic. Stay tuned for the next part.

> References:-
> 1. This blog post is inspired by an awesome talk by [Jason Liu](https://github.com/jxnl)'s, [talk](https://www.youtube.com/watch?v=yj-wSRJwrrc). Please watch it too for better reference.
> 2. [Pydantic](https://pypi.org/project/pydantic/), [Instructor](https://pypi.org/project/instructor/)

