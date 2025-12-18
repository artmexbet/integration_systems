package main

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Nobel Prize API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`

const swaggerSpec = `{
  "openapi": "3.0.0",
  "info": {
    "title": "Nobel Prize API",
    "description": "REST API for managing Nobel Prize laureates and prizes data.\n\n## Authentication\nThis API requires authentication via one of the following methods:\n- **Bearer Token**: Include Authorization: Bearer <token> header\n- **API Key**: Include api_key=<token> query parameter\n\nDefault token for testing: secret-api-token",
    "version": "1.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Local development server"
    }
  ],
  "security": [
    {"BearerAuth": []},
    {"ApiKeyAuth": []}
  ],
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "Token"
      },
      "ApiKeyAuth": {
        "type": "apiKey",
        "in": "query",
        "name": "api_key"
      }
    },
    "schemas": {
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {"type": "string"},
          "message": {"type": "string"}
        }
      },
      "SuccessResponse": {
        "type": "object",
        "properties": {
          "message": {"type": "string"}
        }
      },
      "StatsResponse": {
        "type": "object",
        "properties": {
          "laureates_count": {"type": "integer"},
          "prizes_count": {"type": "integer"},
          "categories_count": {"type": "integer"}
        }
      },
      "LastUpdateResponse": {
        "type": "object",
        "properties": {
          "last_update": {"type": "string", "format": "date-time"}
        }
      },
      "LaureateResponse": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "firstname": {"type": "string"},
          "surname": {"type": "string"},
          "motivation": {"type": "string"},
          "share": {"type": "integer"},
          "updated_at": {"type": "string", "format": "date-time"}
        }
      },
      "LaureateListResponse": {
        "type": "object",
        "properties": {
          "data": {"type": "array", "items": {"$ref": "#/components/schemas/LaureateResponse"}},
          "total": {"type": "integer"},
          "page": {"type": "integer"},
          "per_page": {"type": "integer"},
          "total_pages": {"type": "integer"}
        }
      },
      "CreateLaureateRequest": {
        "type": "object",
        "required": ["id", "firstname", "motivation", "share"],
        "properties": {
          "id": {"type": "integer"},
          "firstname": {"type": "string"},
          "surname": {"type": "string"},
          "motivation": {"type": "string"},
          "share": {"type": "integer", "minimum": 1, "maximum": 4}
        }
      },
      "UpdateLaureateRequest": {
        "type": "object",
        "required": ["firstname", "motivation", "share"],
        "properties": {
          "firstname": {"type": "string"},
          "surname": {"type": "string"},
          "motivation": {"type": "string"},
          "share": {"type": "integer", "minimum": 1, "maximum": 4}
        }
      },
      "PrizeResponse": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "year": {"type": "integer"},
          "category": {"type": "string"},
          "laureates": {"type": "array", "items": {"$ref": "#/components/schemas/LaureateResponse"}},
          "updated_at": {"type": "string", "format": "date-time"}
        }
      },
      "PrizeListResponse": {
        "type": "object",
        "properties": {
          "data": {"type": "array", "items": {"$ref": "#/components/schemas/PrizeResponse"}},
          "total": {"type": "integer"},
          "page": {"type": "integer"},
          "per_page": {"type": "integer"},
          "total_pages": {"type": "integer"}
        }
      },
      "CreatePrizeRequest": {
        "type": "object",
        "required": ["year", "category"],
        "properties": {
          "year": {"type": "integer", "minimum": 1901},
          "category": {"type": "string"},
          "laureate_ids": {"type": "array", "items": {"type": "integer"}}
        }
      },
      "UpdatePrizeRequest": {
        "type": "object",
        "required": ["year", "category"],
        "properties": {
          "year": {"type": "integer", "minimum": 1901},
          "category": {"type": "string"}
        }
      },
      "CategoriesResponse": {
        "type": "object",
        "properties": {
          "categories": {"type": "array", "items": {"type": "string"}}
        }
      }
    }
  },
  "paths": {
    "/health": {
      "get": {
        "tags": ["Health"],
        "summary": "Health check",
        "description": "Returns service health status (no authentication required)",
        "security": [],
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {"type": "string"},
                    "service": {"type": "string"},
                    "version": {"type": "string"}
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/stats": {
      "get": {
        "tags": ["Stats"],
        "summary": "Get dataset statistics",
        "description": "Returns count of laureates, prizes, and categories",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/StatsResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/stats/last-update": {
      "get": {
        "tags": ["Stats"],
        "summary": "Get last update timestamp",
        "description": "Returns the timestamp of the last dataset update",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/LastUpdateResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/categories": {
      "get": {
        "tags": ["Prizes"],
        "summary": "Get all categories",
        "description": "Returns a list of all unique prize categories",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/CategoriesResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/laureates": {
      "get": {
        "tags": ["Laureates"],
        "summary": "List laureates",
        "description": "Returns a paginated list of Nobel laureates",
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "description": "Page number",
            "schema": {"type": "integer", "default": 1}
          },
          {
            "name": "per_page",
            "in": "query",
            "description": "Items per page (max 100)",
            "schema": {"type": "integer", "default": 10, "maximum": 100}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/LaureateListResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      },
      "post": {
        "tags": ["Laureates"],
        "summary": "Create a new laureate",
        "description": "Creates a new Nobel laureate",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreateLaureateRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Laureate created",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/LaureateResponse"}
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/laureates/{id}": {
      "get": {
        "tags": ["Laureates"],
        "summary": "Get laureate by ID",
        "description": "Returns a single laureate by their ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Laureate ID",
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/LaureateResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "404": {
            "description": "Not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      },
      "put": {
        "tags": ["Laureates"],
        "summary": "Update a laureate",
        "description": "Updates an existing Nobel laureate",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Laureate ID",
            "schema": {"type": "integer"}
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/UpdateLaureateRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Laureate updated",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/LaureateResponse"}
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "404": {
            "description": "Not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      },
      "delete": {
        "tags": ["Laureates"],
        "summary": "Delete a laureate",
        "description": "Deletes an existing Nobel laureate",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Laureate ID",
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Laureate deleted",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/SuccessResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/prizes": {
      "get": {
        "tags": ["Prizes"],
        "summary": "List prizes",
        "description": "Returns a paginated list of Nobel prizes",
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "description": "Page number",
            "schema": {"type": "integer", "default": 1}
          },
          {
            "name": "per_page",
            "in": "query",
            "description": "Items per page (max 100)",
            "schema": {"type": "integer", "default": 10, "maximum": 100}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/PrizeListResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      },
      "post": {
        "tags": ["Prizes"],
        "summary": "Create a new prize",
        "description": "Creates a new Nobel prize",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreatePrizeRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Prize created",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/PrizeResponse"}
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/prizes/{id}": {
      "get": {
        "tags": ["Prizes"],
        "summary": "Get prize by ID",
        "description": "Returns a single prize by its ID with associated laureates",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Prize ID",
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/PrizeResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "404": {
            "description": "Not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      },
      "put": {
        "tags": ["Prizes"],
        "summary": "Update a prize",
        "description": "Updates an existing Nobel prize",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Prize ID",
            "schema": {"type": "integer"}
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/UpdatePrizeRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Prize updated",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/PrizeResponse"}
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "404": {
            "description": "Not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      },
      "delete": {
        "tags": ["Prizes"],
        "summary": "Delete a prize",
        "description": "Deletes an existing Nobel prize",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Prize ID",
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Prize deleted",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/SuccessResponse"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/prizes/category/{category}": {
      "get": {
        "tags": ["Prizes"],
        "summary": "Get prizes by category",
        "description": "Returns all prizes for a specific category with their laureates",
        "parameters": [
          {
            "name": "category",
            "in": "path",
            "required": true,
            "description": "Prize category (e.g., physics, chemistry, medicine, literature, peace, economics)",
            "schema": {"type": "string"}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {"$ref": "#/components/schemas/PrizeResponse"}
                }
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    },
    "/api/v1/prizes/year/{year}": {
      "get": {
        "tags": ["Prizes"],
        "summary": "Get prizes by year",
        "description": "Returns all prizes for a specific year",
        "parameters": [
          {
            "name": "year",
            "in": "path",
            "required": true,
            "description": "Prize year",
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {"$ref": "#/components/schemas/PrizeResponse"}
                }
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          }
        }
      }
    }
  }
}`
