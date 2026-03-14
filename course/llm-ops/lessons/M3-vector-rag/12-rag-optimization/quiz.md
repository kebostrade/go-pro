# Quiz: RAG Optimization

## Question 1

What is the primary benefit of hybrid search?

A) It reduces storage costs
B) It combines semantic and keyword matching for better results
C) It makes retrieval faster
D) It eliminates the need for embeddings

## Question 2

What does RRF stand for in hybrid search?

A) Random Result Filtering
B) Reciprocal Rank Fusion
C) Relevant Result Finding
D) Ranked Resource Framework

## Question 3

What is re-ranking in RAG optimization?

A) Re-embedding documents with a different model
B) Post-processing retrieved results with a more accurate scorer
C) Rebuilding the vector index
D) Re-running the query

## Question 4

Which model type is typically used for re-ranking?

A) Dense encoder
B) Sparse encoder
C) Cross-encoder
D) Decoder-only

## Question 5

What is query expansion?

A) Adding more documents to the index
B) Adding related terms to improve recall
C) Increasing the top-K value
D) Expanding the context window

## Question 6

What is HyDE?

A) A vector index type
B) Hypothetical Document Embeddings - using generated answers to improve retrieval
C) A ranking algorithm
D) A chunking strategy

## Question 7

Which metric measures if retrieved context is actually used in the answer?

A) Precision
B) Recall
C) Faithfulness
D) MRR

## Question 8

In the AI Agent Platform, what config option enables hybrid search?

A) HybridMode
B) EnableHybrid
C) UseBM25
D) EnableFusion

## Question 9

What is the typical two-stage retrieval architecture?

A) Embedding + Generation
B) Fast retrieval + Re-ranking
C) Indexing + Searching
D) Chunking + Assembly

## Question 10

What does NDCG stand for?

A) Normalized Document Context Gain
B) Normalized Discounted Cumulative Gain
C) Non-Deterministic Context Generation
D) Networked Document Collection Gateway

## Question 11

What is the purpose of the k parameter in RRF?

A) Number of results to return
B) Constant to prevent high rankings from dominating
C) Vector dimension
D) Chunk size

## Question 12

Which optimization technique uses an LLM to rewrite the query?

A) Query expansion
B) Query rewriting
C) HyDE
D) Re-ranking

## Question 13

What is context precision in RAG evaluation?

A) How fast the retrieval is
B) How relevant the retrieved context is to the question
C) The accuracy of the generated answer
D) The number of documents retrieved

## Question 14

What is the main trade-off with re-ranking?

A) Storage vs. speed
B) Accuracy vs. latency
C) Index size vs. recall
D) Chunk size vs. coherence

## Question 15

Which metric would you use to evaluate if the LLM uses the retrieved context?

A) Answer relevance
B) Context precision
C) MRR
D) Precision@K