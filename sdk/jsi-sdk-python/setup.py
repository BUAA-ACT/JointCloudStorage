from setuptools import find_namespace_packages, setup, find_packages

# 读取项目的readme介绍
with open("README.md", "r") as fh:
    long_description = fh.read()

setup(
    name="jcs-sdk",
    version="1.7.25",
    author="sincerexia",  # 项目作者
    author_email="zhangjh@act.buaa.edu.cn",
    description="This is the official Python SDK for JointCloudStorage",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://www.jointcloudstorage.cn",
    packages=find_packages(),
    python_requires=">=3.6",
    install_requires=[
        "requests>=2.26.0",
    ],
)