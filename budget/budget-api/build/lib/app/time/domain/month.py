class Month:

    def __init__(self, content: int):
        self.content = content

    def __eq__(self, other):
        # Equality Comparison between two objects
        return self.content == other.content

    def __hash__(self):
        # hash(custom_object)
        return hash((self.content))

    @staticmethod
    def JANUARY():
        return Month(1)

    @staticmethod
    def FEBRUARY():
        return Month(2)

    @staticmethod
    def MARCH():
        return Month(3)

    @staticmethod
    def APRIL():
        return Month(4)

    @staticmethod
    def MAY():
        return Month(5)

    @staticmethod
    def JUNE():
        return Month(6)

    @staticmethod
    def JULY():
        return Month(7)

    @staticmethod
    def AUGUST():
        return Month(8)

    @staticmethod
    def SEPTEMBER():
        return Month(9)

    @staticmethod
    def OCTOBER():
        return Month(10)

    @staticmethod
    def NOVEMBER():
        return Month(11)

    @staticmethod
    def DECEMBER():
        return Month(12)
