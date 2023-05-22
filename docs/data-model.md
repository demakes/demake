# Data Model

Klaro has a very simple data model:

* A `Site` defines a website.
* A `Page` defines a single page in a given `Site`.
* A `Route` maps a HTTP route to a given `Page`.
* A `Variable` contains dynamic content that can be used in a `Page`.

That's it! This alone suffices to model any website, including complex patterns like multiple languages and dynamically created pages.

## Site

A site has one or more **domain names**.

How to handle multilingual sites? Keep it simple, a `Site` object won't contain anything that needs to be translated, only pages can have multiple languages.

## Page

A `Page` can have multiple attributes like content, language, title etc. Some attributes are required, other can be optional.

## Route

A `Route` has a regular expression that matches a URL path to a given page.

## Variable

A `Variable` stores content that is used across different pages, e.g. translated title strings. A variable has a unique `name` and can have one or multiple conditions associated with it that are matched against attributes present in the context. This allows us to e.g. define translations for different language versions of a given page.

## Example

```yaml
sites:
	- domain: klaro.org
pages:
	- name: homepage
	  type: template
	  content: \|
		Welcome to Klaro!
variables:
	- name: website.title
	  value: Klaro!
	  conditions:
	  	- language: de

```

## Typed File System Abstraction

We could also model the entire data in Klaro as a file system, which would allow us to use tools like Git to manage it.File systems are familiar to most users and a large part of the CMS content is file-based e.g. static images, texts etc. so people have an intuition of how to work with it. A file system is kind of a directed graph that can have indexes, schemas, links etc. So we should consider simply adopting that.

The file system can modify itself using scripts. Folder structures can define hierarchies and schema definitions inside a folder can define specific data models like posts etc...

We can define symbolic links between files to model relationships, or foreign keys inside file structures.

posts/schema.yml
     /1
       /post.md
       /data.yml
       /files
             /image.png
             /doc.pdf
       /tags
            /software -> /tags/software
            /programming -> /tags/programming

This allows us dynamically define our own schemas and abstractions. We can also map other systems like databases, cloud services etc. to this abstraction and use it to display content.

By keeping the schema definition and the data itself in this structure we can ensure that everything will always be consistent. This also allows to easily e.g. export the entire website into version control and combine external technologies with Klaro CMS, making it easy to adopt the most suitable approach for a given use case. Klaro CMS can act as a backend to frontend apps.

For some things, classical data models are best, for other things like unstructured data, file-like abstractions might work better.

In the end, a website consists of unstructured (e.g. images, downloadable files), lightly-structured (e.g. templates) and strongly structured (e.g. database items) objects that need to be tied together. So we should offer the best tool for all of these things.

## Klaro CMS - Low-Code, Headless, Headfull

In general, Klaro CMS can be seen as a tool to develop apps as well, i.e. we can tie together different mechanisms to allow definition of low-code apps and internal tools with little effort.

## Example: Rendering pages

A given page has a `template` which we fetch from the file system. A page also has metadata like a title, as well as static files like images and a route. Klaro takes the page and associates it with a given route, rendering it using the specified template, which is a symbolic link to a given template file.

This allows us to e.g. specify a file system to render our data from.