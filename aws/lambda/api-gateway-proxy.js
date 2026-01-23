exports.handler = async (event) => {
    console.log("Event received in api-gateway-proxy");

    // Basic validation of expected API Gateway Proxy fields
    const method = event.httpMethod || "UNKNOWN";
    const path = event.path || "/";
    const body = event.body ? JSON.parse(event.body) : {};
    const queryParams = event.queryStringParameters || {};
    const headers = event.headers || {};

    console.log(`Processing ${method} request for ${path}`);

    return {
        statusCode: 200,
        headers: {
            "Content-Type": "application/json",
            "X-Custom-Header": "myna-test"
        },
        body: JSON.stringify({
            message: "API Gateway Proxy response",
            received: {
                method,
                path,
                queryParams,
                bodyHasData: Object.keys(body).length > 0
            }
        }),
        isBase64Encoded: false
    };
};
