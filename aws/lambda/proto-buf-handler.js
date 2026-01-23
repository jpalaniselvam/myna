exports.handler = async (event) => {
    console.log("Event received in proto-buf-handler");

    let buffer;

    // Check if the event has 'data' field (JSON wrapper approach)
    if (event.data) {
        buffer = Buffer.from(event.data, 'base64');
    } else {
        // Fallback or error
        console.log("No data field in event:", event);
        return {
            statusCode: 400,
            body: JSON.stringify({ error: "Missing 'data' field in JSON payload" })
        };
    }

    // Manual Protobuf Decoding for Person
    // message Person {
    //   string name = 1;
    //   int32 id = 2;
    //   string email = 3;
    // }

    const person = {};
    let offset = 0;

    while (offset < buffer.length) {
        const tag = buffer[offset++];
        const fieldNumber = tag >>> 3;
        const wireType = tag & 7;

        if (wireType === 0) { // Varint
            let result = 0;
            let shift = 0;
            while (true) {
                if (offset >= buffer.length) break;
                let b = buffer[offset++];
                result |= (b & 0x7f) << shift;
                if ((b & 0x80) === 0) break;
                shift += 7;
            }
            if (fieldNumber === 2) {
                person.id = result;
            }
        } else if (wireType === 2) { // Length-delimited (string, bytes, etc.)
            let length = 0;
            let shift = 0;
            while (true) {
                if (offset >= buffer.length) break;
                let b = buffer[offset++];
                length |= (b & 0x7f) << shift;
                if ((b & 0x80) === 0) break;
                shift += 7;
            }

            const value = buffer.toString('utf-8', offset, offset + length);
            offset += length;

            if (fieldNumber === 1) {
                person.name = value;
            } else if (fieldNumber === 3) {
                person.email = value;
            }
        } else {
            console.log(`Unknown wire type ${wireType} for field ${fieldNumber}`);
            // Skip or break? For simple decoder, just break to avoid loop
            break;
        }
    }

    console.log("Decoded Person:", person);

    return {
        statusCode: 200,
        headers: { "Content-Type": "application/json" },
        body: person
    };
};
