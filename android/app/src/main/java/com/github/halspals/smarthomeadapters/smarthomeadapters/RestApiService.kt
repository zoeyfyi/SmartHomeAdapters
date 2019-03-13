package com.github.halspals.smarthomeadapters.smarthomeadapters

import com.github.halspals.smarthomeadapters.smarthomeadapters.model.*
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
                    .baseUrl("https://client.api.test.halspals.co.uk/")
                    .build()
                    .create(RestApiService::class.java)
        }

        fun <T> extractErrorFromResponse(response: Response<T>): String? {

            val errorBody = response.errorBody()
            return if (errorBody != null) {
               try {
                    JSONObject(errorBody.string()).getString("error")
                } catch (e: JSONException) {
                    response.message()
                }
            } else {
                response.message()
            }
        }
    }

    @POST("register")
    fun registerUser(@Body user: User): Call<User>

    @POST("login")
    fun loginUser(@Body user: User): Call<Token>

    @GET("robots")
    fun getRobots(@Header("token") token: String): Call<List<Robot>>

    @GET("robot/{id}")
    fun getRobot(@Path("id") id: String, @Header("token") token: String): Call<Robot>

    @POST("robot/{id}")
    fun registerRobot(
            @Path("id") id: String,
            @Header("token") token: String,
            @Body robot: RobotRegistrationBody): Call<ResponseBody>

    @PATCH("robot/{id}/toggle/{current}")
    fun robotToggle(
            @Path("id") id: String, @Path("current") value: Boolean,
            @Header("token") token: String,
            @Body map: Map<String, Boolean>): Call<ResponseBody>

    @GET("usecases")
    fun getAllUseCases(@Header("token") token: String): Call<List<UseCase>>

    @PUT("robot/{robotId}/calibration")
    fun setConfigParameters(
            @Path("robotId") robotId: String,
            @Header("token") token: String,
            @Body params: List<ConfigResult>): Call<ResponseBody>

    @PATCH("robot/{id}/range/{current}")
    fun robotRange(
            @Path("id") id: String,
            @Path("current") value: Int,
            @Header("token") token: String,
            @Body map: Map<String, Boolean>
    ) : Call<ResponseBody>
}