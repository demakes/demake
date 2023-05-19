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
		Welcome to Klaro 
variables:
	- name: website.title
	  value: Klaro!
	  conditions:
	  	- language: de
```