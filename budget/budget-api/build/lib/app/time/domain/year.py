class Year:

    def __init__(self, content:int):
        self.content=content
    
    def __eq__(self, other):
        # Equality Comparison between two objects
        return self.content == other.content

    def __hash__(self):
        # hash(custom_object)
        return hash((self.content))

        