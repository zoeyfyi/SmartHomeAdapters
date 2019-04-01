package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigParameter
import kotlinx.android.synthetic.main.activity_register_robot.*
import kotlinx.android.synthetic.main.fragment_configure_robot.*
import net.openid.appauth.AuthorizationException
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

const val RECALIBRATE_FLAG = "Recalibrate"

/**
 * The final screen of the robot registration wizard, where the user configures the robot parameters.
 */
class ConfigureRobotFragment : Fragment() {

    private val fTag = "ConfigureRobotFragment"

    private val parent by lazy { activity as RegisterRobotActivity }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_configure_robot, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        if (arguments?.getBoolean(RECALIBRATE_FLAG, false) == true) {
            // If we have been instructed to recalibrate, rename the finish button for clarity
            finish_button.setText(R.string.reconfigure_button_text)
        }

        finish_button.setOnClickListener { _ -> setConfigParametersAndFinish() }

        getConfigParameters(view.context)
    }


    /**
     * Gets the configuration parameters for the robot and updates the [parameter_grid].
     *
     * @param context the [Context] to pass to the [parameter_grid]'s adapter
     */
    private fun getConfigParameters(context: Context) {

        progress_bar.visibility = View.VISIBLE
        finish_button.isEnabled = false

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken, _, ex ->
            if (accessToken == null) {
                Log.e(fTag, "[getConfigParameters] got null access token, ex: $ex")
            } else {
                parent.restApiService
                        .getConfigParameters(parent.robotId, accessToken)
                        .enqueue(object: Callback<List<ConfigParameter>> {

                            override fun onResponse(
                                    call: Call<List<ConfigParameter>>,
                                    response: Response<List<ConfigParameter>>) {

                                progress_bar.visibility = View.GONE
                                finish_button.isEnabled = true

                                val params = response.body()

                                if (response.isSuccessful && params != null) {
                                    Log.v(fTag, "[getConfigParameters] Got params $params")
                                    // Set up the grid's adapter to display the configuration
                                    // parameters requested
                                    parameter_grid.adapter = ParameterAdapter(context, params)
                                } else {
                                    val error = RestApiService.extractErrorFromResponse(response)
                                    Log.e(fTag, "[setConfigParameters] got unsuccessful "
                                            + "response or null body; body $params, error: $error")
                                    if (error != null) {
                                        parent.snackbar_layout.snackbar(error)
                                    }
                                }

                            }

                            override fun onFailure(call: Call<List<ConfigParameter>>, t: Throwable) {

                                progress_bar.visibility = View.GONE
                                finish_button.isEnabled = true

                                val error = t.message
                                Log.e(fTag, "[getConfigParameters] FAILED, error: $error")
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }
                            }
                        })
            }
        }
    }

    /**
     * Sets the configuration parameters in the web server and finishes the registration wizard.
     */
    private fun setConfigParametersAndFinish() {

        val config = (parameter_grid.adapter as ParameterAdapter).getConfigValuesSnapshot()

        Log.v(fTag, "[setConfigParametersAndFinish] got config $config")

        progress_bar.visibility = View.VISIBLE
        finish_button.isEnabled = false

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken: String?, _: String?, ex: AuthorizationException? ->
            if (accessToken == null) {
                Log.e(fTag, "[setConfigParametersAndFinish] got null access token, ex: $ex")
            } else {
                parent.restApiService
                        .setConfigParameters(parent.robotId, accessToken, config)
                        .enqueue(object: Callback<ResponseBody> {

                            override fun onResponse(
                                    call: Call<ResponseBody>,
                                    response: Response<ResponseBody>) {

                                progress_bar.visibility = View.GONE
                                finish_button.isEnabled = true

                                if (response.isSuccessful) {
                                    Log.v(fTag, "[setConfigParameters] Success")
                                    parent.toast("Saved successfully")
                                    parent.finish()
                                } else {
                                    val error = RestApiService.extractErrorFromResponse(response)
                                    Log.e(fTag, "[setConfigParameters] got unsuccessful "
                                            + "response, error: $error")
                                    if (error != null) {
                                        parent.snackbar_layout.snackbar(error)
                                    }
                                }
                            }

                            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {

                                progress_bar.visibility = View.GONE
                                finish_button.isEnabled = true

                                val error = t.message
                                Log.e(fTag, "[setConfigParameters] FAILED, error: $error")
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }
                            }
                        })
            }
        }
    }
}
