class PyXBOSError(Exception):
    """Base class for exceptions in pyxbos"""
    pass

class ConfigMissingError(PyXBOSError):
    """Exception raised for errors in the input.

    Attributes:
        expected -- expected key
    """

    def __init__(self, expected, extra=""):
        self.expected = expected
        self.message = "Expected key \"{0}\" in config ({1})".format(expected, extra)

