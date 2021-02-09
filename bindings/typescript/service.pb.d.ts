import * as $protobuf from "protobufjs";
/** Namespace xyzzy. */
export namespace xyzzy {

    /** Represents a XyzzyData */
    class XyzzyData extends $protobuf.rpc.Service {

        /**
         * Constructs a new XyzzyData service.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         */
        constructor(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean);

        /**
         * Creates new XyzzyData service using the specified rpc implementation.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         * @returns RPC service. Useful where requests and/or responses are streamed.
         */
        public static create(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean): XyzzyData;

        /**
         * Calls Ping.
         * @param request PingRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and PingOk
         */
        public ping(request: xyzzy.IPingRequest, callback: xyzzy.XyzzyData.PingCallback): void;

        /**
         * Calls Ping.
         * @param request PingRequest message or plain object
         * @returns Promise
         */
        public ping(request: xyzzy.IPingRequest): Promise<xyzzy.PingOk>;
    }

    namespace XyzzyData {

        /**
         * Callback as used by {@link xyzzy.XyzzyData#ping}.
         * @param error Error, if any
         * @param [response] PingOk
         */
        type PingCallback = (error: (Error|null), response?: xyzzy.PingOk) => void;
    }

    /** Properties of a PingRequest. */
    interface IPingRequest {
    }

    /** Represents a PingRequest. */
    class PingRequest implements IPingRequest {

        /**
         * Constructs a new PingRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: xyzzy.IPingRequest);

        /**
         * Creates a new PingRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns PingRequest instance
         */
        public static create(properties?: xyzzy.IPingRequest): xyzzy.PingRequest;

        /**
         * Encodes the specified PingRequest message. Does not implicitly {@link xyzzy.PingRequest.verify|verify} messages.
         * @param message PingRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: xyzzy.IPingRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a PingRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns PingRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): xyzzy.PingRequest;
    }

    /** Properties of a PingOk. */
    interface IPingOk {
    }

    /** Represents a PingOk. */
    class PingOk implements IPingOk {

        /**
         * Constructs a new PingOk.
         * @param [properties] Properties to set
         */
        constructor(properties?: xyzzy.IPingOk);

        /**
         * Creates a new PingOk instance using the specified properties.
         * @param [properties] Properties to set
         * @returns PingOk instance
         */
        public static create(properties?: xyzzy.IPingOk): xyzzy.PingOk;

        /**
         * Encodes the specified PingOk message. Does not implicitly {@link xyzzy.PingOk.verify|verify} messages.
         * @param message PingOk message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: xyzzy.IPingOk, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a PingOk message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns PingOk
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): xyzzy.PingOk;
    }
}
