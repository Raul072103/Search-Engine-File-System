import json

class TxtFileAnalytics:
    def __init__(self):
        self.no_files = 0
        
        self.min_size = float('inf')
        self.max_size = 0
        self.avg_size = 0
        self.total_size = 0

        self.min_words = float('inf')
        self.max_words = 0
        self.avg_words = 0
        self.total_words = 0
    
    def __str__(self):
        return (
            f"TxtFileAnalytics:\n"
            f"  Number of Files: {self.no_files}\n"
            f"  Size:\n"
            f"    Min: {self.min_size} bytes\n"
            f"    Max: {self.max_size} bytes\n"
            f"    Avg: {self.avg_size:.2f} bytes\n"
            f"    Total: {self.total_size} bytes\n"
            f"  Words:\n"
            f"    Min: {self.min_words}\n"
            f"    Max: {self.max_words}\n"
            f"    Avg: {self.avg_words:.2f}\n"
            f"    Total: {self.total_words}"
        )
    
    def to_dict(self):
        return {
            "no_files": self.no_files,
            "min_size": self.min_size,
            "max_size": self.max_size,
            "avg_size": self.avg_size,
            "total_size": self.total_size,
            "min_words": self.min_words,
            "max_words": self.max_words,
            "avg_words": self.avg_words,
            "total_words": self.total_words
        }

    def to_json(self, indent=4):
        return json.dumps(self.to_dict(), indent=indent)
