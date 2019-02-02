package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.AsyncTask
import android.util.Log
import java.io.InputStream
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

class RESTRequestTask(private val caller: RESTResponseListener)
    : AsyncTask<RESTRequest, Void, Pair<Int, String>>()
{

    private val tag = "RESTRequestTask"

    override fun doInBackground(vararg params: RESTRequest): Pair<Int, String> {
        assert(params.size == 1) { "Expected only one parameter"}
        val request = params[0]
        val endpoint = request.endpoint
        val requestMethod = request.requestMethod
        val data = request.data

        Log.d(tag, "targetUrl: $endpoint")
        Log.d(tag, "requestMethod: $requestMethod")
        Log.d(tag, "data: $data")

        val conn = URL(endpoint).openConnection() as HttpURLConnection
        conn.apply {
            readTimeout = 10000
            connectTimeout = 15000
            this.requestMethod = requestMethod
        }

        // Add the data to the request only if the method permits
        if (requestMethod == RESTRequest.HTTP_POST) {
            val outputStream = conn.outputStream
            val osw = OutputStreamWriter(outputStream)
            osw.write(data)
            osw.flush()
            osw.close()
        }

        // Connect to the database
        conn.connect()

        val responseCode = conn.responseCode
        Log.d(tag, "Return code: $responseCode")

        val response = if (responseCode < HttpURLConnection.HTTP_BAD_REQUEST) {
            extractJSONString(conn.inputStream)
        } else {
            extractJSONString(conn.errorStream)
        }

        Log.d(tag, "JSON Response: $response")
        return Pair(responseCode, response)

    }

    override fun onPostExecute(result: Pair<Int, String>) {
        super.onPostExecute(result)
        caller.handleRESTResponse(
                responseCode = result.first, response = result.second)
    }

    /**
     * WIP: Attempts to extract a single line of JSON from the [InputStream].
     * If there are multiple such lines, or in general lines starting with a {,
     * only the last will be returned.
     * TODO figure out a nicer way of extracting the JSON from the response.
     */
    private fun extractJSONString(inputStream: InputStream): String {
        var res = ""
        inputStream.bufferedReader().use {
            it.forEachLine { line ->
                if (line.startsWith("{")) {
                    res = line
                }
            }
            res
        }

        return res
    }
}