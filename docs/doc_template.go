package docs

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "This is the API server for Unnispick Korean beauty and fashion products",
        "title": "Unnispick K-Style API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "admin@k-stylehub.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {
        "/brands": {
            "get": {
                "produces": ["application/json"],
                "tags": ["brands"],
                "summary": "List all brands",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Items per page",
                        "name": "per_page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Brand"
                        }
                    }
                }
            },
            "post": {
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["brands"],
                "summary": "Create a new brand",
                "parameters": [
                    {
                        "description": "Brand data",
                        "name": "brand",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateBrandRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/Brand"
                        }
                    }
                }
            }
        },
        "/brands/{id}": {
            "get": {
                "produces": ["application/json"],
                "tags": ["brands"],
                "summary": "Get a brand by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Brand ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Brand"
                        }
                    }
                }
            },
            "put": {
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["brands"],
                "summary": "Update a brand",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Brand ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Brand data",
                        "name": "brand",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateBrandRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Brand"
                        }
                    }
                }
            },
            "delete": {
                "produces": ["application/json"],
                "tags": ["brands"],
                "summary": "Delete a brand",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Brand ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/products": {
            "get": {
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "List all products",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Items per page",
                        "name": "per_page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Product"
                        }
                    }
                }
            },
            "post": {
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Create a new product",
                "parameters": [
                    {
                        "description": "Product data",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateProductRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/Product"
                        }
                    }
                }
            }
        },
        "/products/{id}": {
            "get": {
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Get a product by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Product"
                        }
                    }
                }
            },
            "put": {
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Update a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product data",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateProductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Product"
                        }
                    }
                }
            },
            "delete": {
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Delete a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    }
}`
