from fastapi import FastAPI, Request
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer
import torch

app = FastAPI()
model = SentenceTransformer("sentence-transformers/all-MiniLM-L6-v2")  # Fast & small, good for start

class Query(BaseModel):
    text: str

@app.post("/embed")
def embed_text(query: Query):
    embedding = model.encode(query.text, convert_to_numpy=True).tolist()
    return {"embedding": embedding}
