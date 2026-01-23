exports.handler = async (event, context) => {
    console.log("Event received in context-handler");
    console.log("Context:", JSON.stringify(context, null, 2));

    return {
        statusCode: 200,
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            message: "Hello from context-handler.js",
            requestId: context.awsRequestId,
            functionName: context.functionName,
            remainingTime: context.getRemainingTimeInMillis(),
            memoryLimit: context.memoryLimitInMB,
            input: event
        })
    };
};
