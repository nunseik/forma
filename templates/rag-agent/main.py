# Project: {{ .ProjectName }}
# Author: {{ .Author }}
# Created: {{ .Timestamp }}

import numpy as np
import faiss
from sentence_transformers import SentenceTransformer

class SimpleRAG:
    def __init__(self, documents):
        self.documents = documents
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        self.index = self._build_index()

    def _build_index(self):
        print("Building vector index...")
        # Encode the documents into vectors
        embeddings = self.model.encode(self.documents, convert_to_tensor=False)
        
        # Create a FAISS index
        d = embeddings.shape[1]  # Dimension of vectors
        index = faiss.IndexFlatL2(d)
        index.add(np.array(embeddings, dtype='float32'))
        print(f"Index built successfully with {index.ntotal} documents.")
        return index

    def query(self, question, k=1):
        """
        Searches the index for the most relevant document(s).
        """
        print(f"\nSearching for: '{question}'")
        # Encode the query
        query_vector = self.model.encode([question])
        
        # Search the index
        distances, indices = self.index.search(np.array(query_vector, dtype='float32'), k)
        
        # Return the top k results
        results = [self.documents[i] for i in indices[0]]
        return results

def main():
    print(f"--- Starting RAG Agent: {{ .ProjectName }} ---")
    
    # A small corpus of documents for our RAG agent
    corpus = [
        "The capital of France is Paris.",
        "Pygame is a cross-platform set of Python modules designed for writing video games.",
        "The solar system consists of the Sun and the objects that orbit it.",
        "Artificial intelligence is intelligence demonstrated by machines.",
        "A CPU is the electronic circuitry that executes instructions comprising a computer program."
    ]

    rag_agent = SimpleRAG(corpus)

    # Ask some questions
    retrieved_doc = rag_agent.query("What is the main city in France?")
    print(f"Retrieved: {retrieved_doc[0]}")

    retrieved_doc = rag_agent.query("How do you make games in Python?")
    print(f"Retrieved: {retrieved_doc[0]}")

if __name__ == "__main__":
    main()
