WaveBuiltinPSET = b"\x1b\x20\x19\x49\x54\xe8\x6e\xeb\x8f\x91\xff\x98\x3a\xcc\x56\xe6\xc8\x4a\xe2\x9a\x90\x7c\xe7\xe7\x63\x8e\x86\x57\xd5\x14\x99\xb1\x88\xa4"
WaveGlobalNamespace = b"\x1b\x20\xcf\x8d\x19\xd7\x9d\x23\x01\x38\x65\xbe\xf7\x57\xce\xa0\x4c\xde\xe5\xef\x4e\xde\xfc\x80\x8d\xd2\x1e\x4e\x00\x5e\x6f\x80\x47\xcc"

WaveBuiltinE2EE = "decrypt"


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
