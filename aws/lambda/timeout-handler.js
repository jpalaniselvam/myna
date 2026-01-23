exports.handler = async (event) => {
    console.log("Event received in timeout-handler");
    const sleepTime = event.sleepTimeMs || 5000;
    console.log(`Sleeping for ${sleepTime}ms...`);

    await new Promise(resolve => setTimeout(resolve, sleepTime));

    console.log("Wake up!");
    return {
        message: "Finished sleeping",
        sleptFor: sleepTime
    };
};
