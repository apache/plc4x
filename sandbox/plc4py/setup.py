from setuptools import setup, find_packages

setup(
    name="plc4py",
    version="0.1.0",
    description="Plc4py The Python Industrial IOT Adapter",
    long_description="Really, the funniest around.",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "License :: OSI Approved :: Apache 2.0 License",
        "Programming Language :: Python :: 3.8",
        "Topic :: Text Processing :: Linguistic",
    ],
    keywords="modbus plc4x",
    url="https://plc4x.apache.org",
    author='"Apache PLC4X <>"',
    author_email="dev@plc4x.apache.org",
    license="Apache 2.0",
    packages=find_packages(include=["plc4py", "plc4py.*"]),
    install_requires=[
        "pytest-asyncio>=0.18.3",
        "pip-tools",
        "requires",
        "pre-commit>=2.6.0",
        "pytest-mock>=3.3.1",
        "mock>=4.0.2",
        "mypy>=0.942",
        "flake8>=4.0.1",
    ],
)
