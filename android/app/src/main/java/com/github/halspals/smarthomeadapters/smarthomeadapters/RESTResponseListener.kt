package com.github.halspals.smarthomeadapters.smarthomeadapters

interface RESTResponseListener {

    /**
     * Handles a response from a [RESTRequestTask].
     *
     * @param responseCode the response code return by the server
     * @param response the json found in the body of the response
     */
    fun handleRESTResponse(responseCode: Int, response: String)
}