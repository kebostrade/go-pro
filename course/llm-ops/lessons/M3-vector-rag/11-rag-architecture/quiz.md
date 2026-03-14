# Quiz: RAG Architecture

## Question 1

What does RAG stand for?

A) Retrieval-Augmented Generation
B) Random Access Gateway
C) Re-ranker Advanced Generator
D) Real-time API Gateway

## Question 2

What is the primary purpose of chunking in a RAG system?

A) To compress document size
B) To split documents into smaller, searchable units
C) To remove irrelevant information
D) To encrypt document content

## Question 3

Which chunking strategy works best for well-structured documents with clear headers?

A) Fixed-size chunking
B) Semantic chunking
C) Random chunking
D) Character-level chunking

## Question 4

In the AI Agent Platform, which component handles document indexing?

A) VectorStore
B) DocumentIndexer
C) Embedder
D) Agent

## Question 5

What is the purpose of metadata filtering in RAG retrieval?

A) To reduce vector search latency
B) To filter results based on document attributes like source, date, or category
C) To compress retrieval results
D) To rank results by importance

## Question 6

What is context assembly in a RAG pipeline?

A) Storing documents in the vector database
B) Combining retrieved documents into a prompt for the LLM
C) Generating embeddings for queries
D) Indexing new documents

## Question 7

What happens when retrieved content exceeds the LLM's context window?

A) The system returns an error
B) The LLM generates random content
C) Content must be truncated or summarized to fit
D) The query is rejected

## Question 8

Which of the following is a valid chunking overlap strategy?

A) Remove overlapping content completely
B) Include 50 tokens from previous chunk in current chunk
C) Only use first 100 tokens of each chunk
D) Combine all chunks into one large document

## Question 9

What is re-ranking in the context of RAG?

A) Re-embedding documents with a different model
B) Post-processing retrieved results to improve relevance ordering
C) Re-indexing documents in the vector store
D) Regenerating chunks for better results

## Question 10

In the AI Agent Platform, which type of RAG pipeline is used for code search?

A) DocumentRAGPipeline
B) CodeRAGPipeline
C) TextRAGPipeline
D) SemanticRAGPipeline

## Question 11

What is the typical chunk size for code repositories?

A) 32-64 tokens
B) 128-256 tokens
C) 256-512 tokens
D) 2048+ tokens

## Question 12

What is the purpose of the overlap parameter in chunking?

A) To ensure semantic continuity between chunks
B) To increase storage requirements
C) To reduce retrieval speed
D) To compress embeddings

## Question 13

Which tool is commonly used for extracting text from PDFs?

A) BeautifulSoup
B) pdfplumber
C) Markdown parser
D) JSON loader

## Question 14

What is the recommended chunk size for FAQ-style content?

A) 1024+ tokens
B) 512-768 tokens
C) 128-256 tokens
D) 4096 tokens