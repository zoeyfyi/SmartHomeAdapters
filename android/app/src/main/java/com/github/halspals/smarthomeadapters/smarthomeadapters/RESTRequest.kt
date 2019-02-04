package com.github.halspals.smarthomeadapters.smarthomeadapters

import org.json.JSONObject

sealed class RESTRequest() {

    protected val baseURL = "http://client.api.halspals.co.uk"
    abstract val endpoint: String
    abstract val requestMethod: String
    abstract val data: String
    abstract val type: String

    // TODO if we need more constants in the future
    // we should consider moving these to a separate file
    companion object {
        const val HTTP_POST = "POST"
        const val LOGIN_TYPE = "LoginREST"
        const val REGISTER_TYPE = "RegisterREST"
    }

    class LOGIN(private val email: String, private val password: String) : RESTRequest() {

        override val type = LOGIN_TYPE
        override val endpoint = "$baseURL/login"
        override val requestMethod = HTTP_POST  // TODO make the endpoint accept this
        override val data: String
            get() {
                val json = JSONObject()
                json.put("email", email)
                json.put("password", password)

                return json.toString()
            }


    }

    class REGISTER(private val email: String, private val password: String) : RESTRequest() {

        override val type = REGISTER_TYPE
        override val endpoint = "$baseURL/register"
        override val requestMethod = HTTP_POST
        override val data: String
            get() {
                val json = JSONObject()
                json.put("email", email)
                json.put("password", password)

                return json.toString()
            }
    }
}