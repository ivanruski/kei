# Overview

KEI stadns for `kubectl explain interactive <type>.<fieldName>[.<fieldName>]`. 

This program is trying to mimic the commands available in [less](http://www.greenwoodsoftware.com/less). So people familiar with it should find working with KEI intuitive.

`kei` executes the `kubectl explain` command in the background. It is not not looking at the KUBECONIFIG file.

## Usage

```text
type  /<type>[.<fieldName>]<Enter> to see an explaination of a Kubernetes resource
type  u                            to scroll down one half of the screen size
type  u                            to scroll up one half of the screen size
type  q<Enter>                     to exit
type  h                            to see this message again
```
