exports.handler = async (event) => {
    console.log("Event received in json-handler");
    return {
        statusCode: 200,
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            message: "Hello from json-handler.js",
            input: event
        })
    };
};
