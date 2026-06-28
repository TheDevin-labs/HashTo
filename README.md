# HashTo

HashTo is a high-performance command-line utility written in Go designed to compile raw JSON data into dynamic, well-formatted JavaScript source code. By parsing abstract data structures and mapping them directly to native JavaScript representations, HashTo bridges the gap between static configurations and dynamic front-end environments with zero overhead.

[![License](https://img.shields.io/badge/License-BSD__3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
![Language](https://img.shields.io/badge/Language-Go-blue)
![Automation](https://img.shields.io/badge/Build-Makefile-orange)


---

## Architecture Overview

The utility processes data via a decoupled compiler architecture, ensuring deep recursive nesting stability and type safety during the transpilation process.


```
[ Input JSON File ] ---> [ Go Unmarshaling Engine ] ---> [ AST Type Validation ] ---> [ JS Format Generator ] ---> [ Output .js Script ]
```
---

## Notice

> [!IMPORTANT]
> **COMPLIANCE & SYNTAX NOTICE**
> 
> HashTo automatically evaluates object keys during the compilation sequence. If a JSON key contains special characters, mathematical operators, or hyphens (e.g., `"system-status"`), the compiler wraps the key in quotes to preserve syntactical validity in JavaScript. Clean, alphanumeric keys will remain unquoted to optimise performance and maintain idiomatic JavaScript standards.
> 
> Furthermore, all numeric formats are standardised from the initial JSON 64-bit floating-point notation (`float64`) down to their cleanest respective integer or decimal equivalents in the final asset output.

---

## System Requirements

* Go Compiler (version 1.18 or higher)
* GNU Make (optional, for automation via Makefile)

---

## Installation and Compilation

To compile the binary optimised for your local operating system and processor architecture, execute the automated build sequence:

```bash
make build

```
This generates a platform-specific binary following the naming convention: hashto-<os>-<arch>. For instance, on a standard British Linux configuration, this yields hashto-linux-amd64.
## Command Line Interface Usage
The utility requires a specific parameter structure to parse metadata and map output variables accurately.
### Syntax
```bash
./hashto js=<variable_name> @file=<input_file.json> <output_file.js>

```
### Argument Parameters
 * js=: Defines the name of the constant variable exported within the generated JavaScript file.
 * @file=: Specifies the relative or absolute path to the source JSON file.
 * Target Argument: The final trailing argument dictates the destination path of the generated .js script.
### Operational Example
Assuming you have a source configuration file named work.json:
```json
{
  "appName": "HashTo Engine",
  "version": 2.4,
  "isActive": true,
  "system-status": "optimal"
}

```
Execute the following command in your terminal terminal environment:
```bash
./hashto js=hash @file=work.json work.js

```
### Resulting Output
The transpiled work.js file will contain clean, modularised JavaScript:
```javascript
const hash = {
  appName: "HashTo Engine",
  version: 2.4,
  isActive: true,
  "system-status": "optimal"
};

export default hash;

```
## Technical Specifications & Features
 * **Strict Recursive Resolution:** Traverses multi-layered, nested JSON structures (objects within arrays within objects) without stack overflow risks.
 * **Deterministic Object Key Sanitisation:** Automatically flags invalid ECMAScript identifiers and wraps them in safe lexical enclosures.
 * **Cross-Platform Portability:** Zero external runtime dependencies; compiles directly down to a single binary file.
 * 
