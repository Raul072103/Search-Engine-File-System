# Changelog

All notable changes to this project will be documented in this file.

## [0.0.1] - 2025-05-26

### Added
- Initial version of app with:
    - **File Indexer**
      - file indexing based on config.json
      - automatic progress saving
      - automatic file changes detection under the root folder using USN Journal

    - **backend API**
        - search endpoint, allowing for queries based on content, file extensions and names
        - query suggestions endpoint, suggests using qdrant based on previous queries that were 
introduced by the user
        - basic grammar spelling corrections endpoint
    
    - **Web application**
      - basic React application, with one page
      - search query summary
      - contextual widgets (non-functional)
      - basic parser for queries like: "content:milk name:recipe extensions:.txt"
  

---

