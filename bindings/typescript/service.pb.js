/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

$root.xyzzy = (function() {

    /**
     * Namespace xyzzy.
     * @exports xyzzy
     * @namespace
     */
    var xyzzy = {};

    xyzzy.XyzzyData = (function() {

        /**
         * Constructs a new XyzzyData service.
         * @memberof xyzzy
         * @classdesc Represents a XyzzyData
         * @extends $protobuf.rpc.Service
         * @constructor
         * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
         * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
         * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
         */
        function XyzzyData(rpcImpl, requestDelimited, responseDelimited) {
            $protobuf.rpc.Service.call(this, rpcImpl, requestDelimited, responseDelimited);
        }

        (XyzzyData.prototype = Object.create($protobuf.rpc.Service.prototype)).constructor = XyzzyData;

        /**
         * Creates new XyzzyData service using the specified rpc implementation.
         * @function create
         * @memberof xyzzy.XyzzyData
         * @static
         * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
         * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
         * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
         * @returns {XyzzyData} RPC service. Useful where requests and/or responses are streamed.
         */
        XyzzyData.create = function create(rpcImpl, requestDelimited, responseDelimited) {
            return new this(rpcImpl, requestDelimited, responseDelimited);
        };

        /**
         * Callback as used by {@link xyzzy.XyzzyData#ping}.
         * @memberof xyzzy.XyzzyData
         * @typedef PingCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {xyzzy.PingOk} [response] PingOk
         */

        /**
         * Calls Ping.
         * @function ping
         * @memberof xyzzy.XyzzyData
         * @instance
         * @param {xyzzy.IPingRequest} request PingRequest message or plain object
         * @param {xyzzy.XyzzyData.PingCallback} callback Node-style callback called with the error, if any, and PingOk
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(XyzzyData.prototype.ping = function ping(request, callback) {
            return this.rpcCall(ping, $root.xyzzy.PingRequest, $root.xyzzy.PingOk, request, callback);
        }, "name", { value: "Ping" });

        /**
         * Calls Ping.
         * @function ping
         * @memberof xyzzy.XyzzyData
         * @instance
         * @param {xyzzy.IPingRequest} request PingRequest message or plain object
         * @returns {Promise<xyzzy.PingOk>} Promise
         * @variation 2
         */

        return XyzzyData;
    })();

    xyzzy.PingRequest = (function() {

        /**
         * Properties of a PingRequest.
         * @memberof xyzzy
         * @interface IPingRequest
         */

        /**
         * Constructs a new PingRequest.
         * @memberof xyzzy
         * @classdesc Represents a PingRequest.
         * @implements IPingRequest
         * @constructor
         * @param {xyzzy.IPingRequest=} [properties] Properties to set
         */
        function PingRequest(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Creates a new PingRequest instance using the specified properties.
         * @function create
         * @memberof xyzzy.PingRequest
         * @static
         * @param {xyzzy.IPingRequest=} [properties] Properties to set
         * @returns {xyzzy.PingRequest} PingRequest instance
         */
        PingRequest.create = function create(properties) {
            return new PingRequest(properties);
        };

        /**
         * Encodes the specified PingRequest message. Does not implicitly {@link xyzzy.PingRequest.verify|verify} messages.
         * @function encode
         * @memberof xyzzy.PingRequest
         * @static
         * @param {xyzzy.IPingRequest} message PingRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        PingRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            return writer;
        };

        /**
         * Decodes a PingRequest message from the specified reader or buffer.
         * @function decode
         * @memberof xyzzy.PingRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {xyzzy.PingRequest} PingRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        PingRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.xyzzy.PingRequest();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        return PingRequest;
    })();

    xyzzy.PingOk = (function() {

        /**
         * Properties of a PingOk.
         * @memberof xyzzy
         * @interface IPingOk
         */

        /**
         * Constructs a new PingOk.
         * @memberof xyzzy
         * @classdesc Represents a PingOk.
         * @implements IPingOk
         * @constructor
         * @param {xyzzy.IPingOk=} [properties] Properties to set
         */
        function PingOk(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Creates a new PingOk instance using the specified properties.
         * @function create
         * @memberof xyzzy.PingOk
         * @static
         * @param {xyzzy.IPingOk=} [properties] Properties to set
         * @returns {xyzzy.PingOk} PingOk instance
         */
        PingOk.create = function create(properties) {
            return new PingOk(properties);
        };

        /**
         * Encodes the specified PingOk message. Does not implicitly {@link xyzzy.PingOk.verify|verify} messages.
         * @function encode
         * @memberof xyzzy.PingOk
         * @static
         * @param {xyzzy.IPingOk} message PingOk message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        PingOk.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            return writer;
        };

        /**
         * Decodes a PingOk message from the specified reader or buffer.
         * @function decode
         * @memberof xyzzy.PingOk
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {xyzzy.PingOk} PingOk
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        PingOk.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.xyzzy.PingOk();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        return PingOk;
    })();

    return xyzzy;
})();

module.exports = $root;
