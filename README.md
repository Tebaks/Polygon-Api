[![Unit Tests](https://github.com/Tebaks/Polygon-Api/actions/workflows/unit-test.yaml/badge.svg)](https://github.com/Tebaks/Polygon-Api/actions/workflows/unit-tests.yaml)
[![Functional Tests](https://github.com/Tebaks/Polygon-Api/actions/workflows/functional-test.yaml/badge.svg)](https://github.com/Tebaks/Polygon-Api/actions/workflows/functional-tests.yaml)

---
# Polygon API
This is a simple API for create and get polygons with given vertices.

# Table of Contents
- [How to run?](#how-to-run)
- [Tests](#tests)
  * [Unit Tests](#unit-tests)
  * [Functional Tests](#functional-tests)
    + [Generate 1000 Polygons](#generate-1000-polygons)
- [Calculating The Area Of Polygon](#calculating-the-area-of-polygon)
- [Endpoints](#endpoints)
  * [Swagger](#swagger)
  * [Find Polygon By Name](#find-polygon-by-name)
  * [Create Polygon](#create-polygon)

---
## [How to run?](#table-of-contents)
To run api you need to edit config.yaml file:
```yaml
Server:
  Port: 8080
Mongo:
  ConnectionString: mongodb://admin:test@mongodb:27017
```
Or set environment variables:
```bash
export MONGO_CONNECTIONSTRING=mongodb://admin:test@mongodb:27017
export SERVER_PORT=8080
```
After that you can run api with command:
```bash
docker-compose up --build -d
```

---
## [Tests](#table-of-contents)
I write both unit and functional tests.

### [Unit Tests](#table-of-contents)
For unit test I mock repository, service, logger and mongoDB. To run unit tests you need to run command:
```bash
go test -tags=unit -v -cover ./...
```
On image below you can see results and coverage.
</br>
<kbd>![unit tests](./.github/unit_tests.png)</kbd>

### [Functional tests](#table-of-contents)
For functional service I run mock mongoDB using [Testcontainers](https://www.testcontainers.org/). To run functional tests you need to run docker daemon and run command:
```bash
go test -tags=functional -v ./...
```
On image below you can see results.
<kbd>![functional tests](./.github/functional_tests.png)</kbd>
#### [Generate 1000 Polygons](#table-of-contents)
When running functional tests application will generate 1000 polygons and prints 3 random polygons.
</br>
<kbd>![generate 1000 polygons](./.github/random_polygon.png)</kbd>

---
## [Calculating The Area Of Polygon](#table-of-contents)
I use [geo-algorithm](https://en.wikipedia.org/wiki/Polygon#Area) to calculate area of polygon. But in this formula vertices must be in clockwise or counterclockwise order. So at first I [sort](https://www.baeldung.com/cs/sort-points-clockwise) vertices in clockwise order. For example if you have vertices:
```json
{
  "vertices": [
    {
      "x": 1,
      "y": 5
    },
    {
      "x": 0,
      "y": 12
    },
    {
      "x": 4,
      "y": 6
    },
    {
      "x": 3,
      "y": 5
    },
    {
      "x": 2,
      "y": 2
    }
  ]
}
```
Then you will get:
```json
{
  "success": true,
  "error": "",
  "data": {
    "name": "Polygon_5_z0vfu",
    "vertices": [
      {
        "x": 3,
        "y": 5
      },
      {
        "x": 2,
        "y": 2
      },
      {
        "x": 1,
        "y": 5
      },
      {
        "x": 0,
        "y": 12
      },
      {
        "x": 4,
        "y": 6
      }
    ],
    "area": 15
  }
}
```
**On the picture below you can see shape and area.**
</br>
<kbd>![picture](./.github/polygon.png)</kbd>

---
## [Endpoints](#table-of-contents)
### [Swagger](#table-of-contents)
You can use swagger to test API. for example for port `8080` swagger url is: `http://localhost:8080/swagger/index.html#/`
<kbd>![swagger](./.github/swagger.png)</kbd>

### [Find Polygon By Name](#table-of-contents)
- **method: `GET`**
- **url: `/polygon/{name}`**
- **description: `Get polygon by name`**
- **response:**
  ```json
    {
      "success": boolean,
      "error": string,
      "data": {
        "name": string,
        "vertices": [
          {
            "x": number,
            "y": number
          }
        ],
        "area": number
      }
    }

### [Create Polygon](#table-of-contents)
- **method: `POST`**
- **url: `/polygon/`**
- **description: `Create polygon with vertices`**
- **request:**
  ```json
    {
      "vertices": [
        {
          "x": number,
          "y": number
        }
      ]
    }
  ```
- **response:**
  ```json
    {
      "success": boolean,
      "error": string,
      "data": {
        "name": string,
        "vertices": [
          {
            "x": number,
            "y": number
          }
        ],
        "area": number
      }
    }
  ```
