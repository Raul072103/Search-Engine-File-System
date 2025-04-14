from pathlib import Path
import PyPDF2
import docx
import os
import json

def traverse(dirPath: str | Path, depth: int, analytics: dict):
    if isinstance(dirPath, str):
        dirPath = Path(dirPath)

    try:
        for elt in dirPath.iterdir():
            if elt.is_file():    
                file_name, extension = os.path.splitext(elt)

                analytics["no_files"] += 1
            
                # extensions
                analytics["extensions"][extension] += 1

                if extension == ".pdf":
                    try:
                        with open(elt, "rb") as file:
                            pdf_reader = PyPDF2.PdfReader(file)
                            text = ""
                            for page in pdf_reader.pages:
                                text += page.extract_text() or ""
                            words = len(text.split())
                            size = elt.stat().st_size
                            update_analytics(analytics[".pdf"], words, size)
                    except Exception as e:
                        print(f"Error processing PDF {elt}: {e}")
                elif extension == ".docx":
                    try:
                        doc = docx.Document(elt)
                        text = ""
                        for paragraph in doc.paragraphs:
                            text += paragraph.text + " "
                        words = len(text.split())
                        size = elt.stat().st_size
                        update_analytics(analytics[".docx"], words, size)
                    except Exception as e:
                        print(f"Error processing DOCX {elt}: {e}")
                elif extension == ".txt":
                    try:
                        with open(elt, "r", encoding="utf-8") as f:
                            chars = f.read()
                            words = len(chars.split())
                            size = elt.stat().st_size
                            update_analytics(analytics[".txt"], words, size)
                    except Exception as e:
                        print(f"Error processing TXT {elt}: {e}")
                    
                try:
                    size = elt.stat().st_size
                    analytics["total_size"] += size
                except PermissionError:
                    print(f"PermissionError accessing size of {elt}") #or pass, or log
                    pass

                analytics["depth"]["depths"].append(depth)

            elif elt.is_dir():
                depth += 1
                if analytics["depth"]["max"] < depth:
                    analytics["depth"]["max"] = depth
                    
                analytics["no_folders"] += 1

                traverse(elt, depth, analytics)
    except PermissionError as e:
        print(f"Permission error accessing directory {dirPath}: {e}")
    except Exception as e:
        print(f"An unexpected error occured while traversing {dirPath}: {e}")    

def update_analytics(analytics_obj, words, size):
    analytics_obj.no_files += 1
    if size < analytics_obj.min_size:
        analytics_obj.min_size = size
    if size > analytics_obj.max_size:
        analytics_obj.max_size = size
    analytics_obj.total_size += size
    if words < analytics_obj.min_words:
        analytics_obj.min_words = words
    if words > analytics_obj.max_words:
        analytics_obj.max_words = words
    analytics_obj.total_words += words