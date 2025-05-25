# Search Engine for Windows File System

This was the project that I had to construct for Software Design course as a 
laboratory work that spanned the whole semester.

It turned out a quite nice project, with some interesting features.

# Context

The main goal was to design a responsive file search system, by indexing the
file system before hand.

It turned out faster than I expected, by using Postgres as a database for storing
the structures and data required for indexing the files.

# Interesting Features

- the project supports searching through the file system by words
- Designed quite a nice architecture for processing the indexing, using
events (check **ARCHITECTURE.md**)
- I integrated the project so it can parse **USN Journal** logs 
obtained from querying Windows Internals. This lead to the indexer
being capable to detect changes at runtime, even if the indexer
was shut down before.
- I used **Qdrant** that allowed to **suggest queries** to
the user of the web application by storing all the previous
queries and calculating the cosine similarity of the current
query and its previous queries.
- Created an **in-memory cache** in Go for repeating queries