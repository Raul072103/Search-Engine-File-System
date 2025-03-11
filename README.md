# Search Engine for my File System

This is going to be a project built by me consisting in a search engine designed for my file system.

## Features

- [ ] Search by filename
- [ ] Search by metadata
- [ ] Search by words inside files

## Notes

- USN Journal
- Master-slave format
- Not going to take into consideration links from the start
- What about hidden files?

## Search Criteria
- Single/multi world 
- Name
- Extension

## System Requirements

### 1. Text search

#### Option 1
For a functional single/multi word search I thought of using a table like:

| Id | Word    | Document|
|----|---------| ---     |
| 1  | "ceapa" | "reteta.txt"|

Total words = approx **68.045.814**

Assume avg word = approx **5 chars**

Assume I use a foreign key to a document Id

For a row in my database (not taking into consideration the index tables created)

**1 row** = approx (8 bytes PK) + (2 bytes * 5) + (8 bytes FK)
**1 row** = 40 bytes (I added 12 bytes extra just to be sure, plus I don't take into cosnideration a lot of things)

**Total storage space just for the search feature** = 40 bytes * 68.045.814

**space** = 2,721,832,560 **bytes**

**space** = 2.71 **GB**

This would mean **at least** 

#### Option 2

Use GIN indexes from Postgres, which are basically what I described above.

GIN (Generalized Inverted Index)

#### Option 3

Use **tsvector**, **tsquery**. 

#### Option 2 + Option 3

What I read is something where I combine a GIN index and a **tsvector** but this will
require me to **store** also **the content** of the files which seems to
memory expensive.

```
CREATE INDEX pgweb_idx ON pgweb USING GIN (to_tsvector('english', body));
```

### Conclusion

I think I will stick with **Option 1** where I create a table I mentioned and where
I will put a **GIN** index.


### Questions

- Opinions on my battle plan
- **Do we want exact matches for what the user is looking for?**
- Can I use Elastic Search or should I try and implement something simillar by myself?
- Access denied files
- Different file encodings
