package com.github.halspals.smarthomeadapters.smarthomeadapters

import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Token
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import okhttp3.ResponseBody
import org.json.JSONException
import org.json.JSONObject
import retrofit2.Call
import retrofit2.Response
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.*

interface RestApiService {

    companion object {
        fun new(): RestApiService {
            return Retrofit.Builder()
                    .addConverterFactory(GsonConverterFactory.create())
                    .baseUrl("http://client.api.halspals.co.uk/")
                    .build()
                    .create(RestApiService::class.java)
        }

        fun <T> extractErrorFromResponse(response: Response<T>): String? = try {
            JSONObject(response.errorBody()?.string()).getString("error")
        } catch (e: JSONException) {
            response.message()
        }
    }

    @POST("register")
    fun registerUser(@Body user: User): Call<User>

    @POST("login")
    fun loginUser(@Body user: User): Call<Token>

    @GET("robots")
    fun getRobots(): Call<List<Robot>>

    @GET("robot/{id}")
    fun getRobot(@Path("id") id: String): Call<Robot>

    @PATCH("robot/{id}/toggle/{current}")
    fun robotToggle(@Path("id") id: String, @Path("current") value: Boolean, @Body map: Map<String, Boolean>): Call<ResponseBody>

}