import setuptools

#with open("README.md", "r") as fh:
#    long_description = fh.read()

print(setuptools.find_packages())
setuptools.setup(
    name="pyxbos",
    version="0.0.1",
    author="Gabe Fierro",
    author_email="gtfierro@cs.berkeley.edu",
    description="Python3 WAVEMQ XBOS",
    #long_description=long_description,
    #long_description_content_type="text/markdown",
    url="https://github.com/xboswave",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
    ],
)
