exports.handler = async (event) => {
    console.log("Event received in error-handler");
    const errorType = event.errorType || "generic";

    if (errorType === "timeout") {
        // user should use timeout-handler.js for real timeouts, 
        // but this simulates a logical timeout error if needed
        throw new Error("Simulated timeout error");
    } else if (errorType === "custom") {
        throw new Error(event.errorMessage || "Custom simulated error");
    }

    // Default failure
    throw new Error("Standard simulated execution failure");
};
