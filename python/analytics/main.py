from traverse import traverse
from files import TxtFileAnalytics
from collections import defaultdict
import json

extensions_dict = {}
txt_extensions = set([".txt", ".docx", ".pdf"])

statistics_dict = {
    "no_files": 0,
    "no_folders": 0,
    "extensions": defaultdict(int),
    "depth": {
        "max": 0,
        "depths": [],
        "avg": 0
    },
    ".txt": TxtFileAnalytics(),
    ".docx": TxtFileAnalytics(),
    ".pdf": TxtFileAnalytics(),
    "total_words": 0,
    "total_txt_size": 0,
    "total_size": 0
}

if __name__ == '__main__':
    traverse("C:\\", 0, statistics_dict)
    statistics_dict["total_words"] = statistics_dict[".txt"].total_words + statistics_dict[".docx"].total_words +  statistics_dict[".pdf"].total_words
    statistics_dict["total_txt_size"] = statistics_dict[".txt"].total_size + statistics_dict[".docx"].total_size +  statistics_dict[".pdf"].total_size
    statistics_dict["depth"]["avg"] = sum(statistics_dict["depth"]["depths"]) / statistics_dict["no_folders"]

    statistics_dict["depth"]["depths"] = []

    serializable_dict = {
        key: (value.to_dict() if isinstance(value, TxtFileAnalytics) else value)
        for key, value in statistics_dict.items()
    }

    # Serialize the dictionary
    json_output = json.dumps(serializable_dict, indent=4)
    print(json_output)


