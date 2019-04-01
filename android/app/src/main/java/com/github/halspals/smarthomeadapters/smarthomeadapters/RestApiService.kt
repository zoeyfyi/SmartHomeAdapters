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

/**
 * Defines the Rest API interface for use with [Retrofit].
 */
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

    @GET("robots")
    fun getRobots(@Header("token") token: String): Call<List<Robot>>

    @GET("robot/{id}")
    fun getRobot(@Path("id") id: String, @Header("token") token: String): Call<Robot>

    @POST("robot/{id}")
    fun registerRobot(
            @Path("id") id: String,
            @Header("token") token: String,
            @Body robot: RobotRegistrationBody
    ): Call<ResponseBody>

    @PATCH("robot/{id}/toggle/{current}")
    fun robotToggle(
            @Path("id") id: String, @Path("current") value: Boolean,
            @Header("token") token: String,
            @Body map: Map<String, Boolean>
    ): Call<ResponseBody>

    @GET("usecases")
    fun getAllUseCases(@Header("token") token: String): Call<List<UseCase>>

    @PUT("robot/{robotId}/calibration")
    fun setConfigParameters(
            @Path("robotId") robotId: String,
            @Header("token") token: String,
            @Body params: List<ConfigResult>
    ): Call<ResponseBody>

    @GET("robot/{robotId}/calibration")
    fun getConfigParameters(
            @Path("robotId") robotId: String,
            @Header("token") token: String
    ): Call<List<ConfigParameter>>

    @PATCH("robot/{id}/range/{current}")
    fun robotRange(
            @Path("id") id: String,
            @Path("current") value: Int,
            @Header("token") token: String,
            @Body map: Map<String, Boolean>
    ): Call<ResponseBody>

    @DELETE("robot/{id}")
    fun deleteRobot(
            @Path("id") id: String,
            @Header("token") token: String
    ): Call<ResponseBody>

    @PATCH("robot/{id}/nickname")
    fun renameRobot(
            @Path("id") id:  String,
            @Header("token") token: String,
            @Body nameMap: Map<String, String>
    ): Call<ResponseBody>

    @PATCH("robot/{id}/reconfigure")
    fun patchUseCase(
            @Path("id") id:  String,
            @Body useCaseMap: Map<String, String>
    ): Call<ResponseBody>

    @GET("user")
    fun getUserName(@Header("token") token: String): Call<User>
}