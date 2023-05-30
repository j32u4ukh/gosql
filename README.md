# gosql

[![Build and Test](https://github.com/j32u4ukh/gosql/actions/workflows/main.yml/badge.svg)](https://github.com/j32u4ukh/gosql/actions/workflows/main.yml)
[![Build and Test](https://github.com/j32u4ukh/gosql/actions/workflows/develop.yml/badge.svg)](https://github.com/j32u4ukh/gosql/actions/workflows/develop.yml)
[![Build and Test](https://github.com/j32u4ukh/gosql/actions/workflows/test.yml/badge.svg)](https://github.com/j32u4ukh/gosql/actions/workflows/test.yml)

# 專案結構

### Layer 0

不引用其他套件，被其他套件所引用的套件

* database
* utils

### Layer 1

只引用 Layer 0 的套件。

* stmt

### Layer 2

只引用 Layer 0 或 Layer 1 的套件。

* plugin
* sync