Here’s the revised content in English:

# NetAres
___
## Content
___
- [Content](##Content)
- [Usage](##Usage)
  - [General](###General)
  - [Flags](###Flags)
  - [Masks](###Masks)
- [Authors](##Authors)
- [License](##License)
___
## Usage
___
### General
NetAres is developed as an open-source product for bash shells. It is recommended to use byte concatenation:
*`netares <flags...> > log.out`* 
instead of piping the output:
*`netares <flags...> | log.out`*
___
### Flags
Description of flags:
| Format        | Default               | Description                             |
|---------------|-----------------------|-----------------------------------------|
| `--mask`      | `"./..."`             | Path to mask file                       |
| `--type`      | `"raw"`               | Type of output                          |
| `--target`    | `"username"`          | Target name                             |
| `--timeout`   | `1000`                | Timeout for HTTP in milliseconds        |
| `--retries`   | `3`                   | Number of retries                       |
___
### Masks
Masks are the key feature of the software. They are easy to produce, deploy, and obtain results from queries, and they are very powerful in routing, as they are designed to work with **XPath**. Here’s a simple default example:

- *`github_mask.json`*: 
```json
{
    "name": "github",
    "source": "https://github.com/*",
    "data": {
      "username": {
        "route": "//span[@class='p-nickname vcard-username d-block']/text()"
      },
      "repositories_count": {
        "route": "//a[contains(@class, 'UnderlineNav-item')]/span[@class='Counter']"
      }
    }
}
```

As you can see, some of the fields are self-named, including all routes mentioned in **data**. This *adds flexibility* to the output.
___
## Authors
There is only one person who works on creating this resource - me, @nemesidaa. You can see my profile [here](https://github.com/nemesidaa).
___
## License
___
This project is licensed under the MIT License.  
**I would be glad to see your forks and connect them with my code to achieve better performance.**

---

**EOF**, glad to see you! Your **Dockie*