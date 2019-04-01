package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.RobotRegistrationBody
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.UseCase
import kotlinx.android.synthetic.main.activity_register_robot.*
import kotlinx.android.synthetic.main.fragment_select_attach.*
import net.openid.appauth.AuthorizationException
import okhttp3.ResponseBody
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.intentFor
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * The third screen in the robot registration wizard, allowing the user to choose a use case for the robot.
 */
class SelectAttachFragment : Fragment() {


    private val fTag = "SelectAttachFragment"

    private val parent by lazy { activity as RegisterRobotActivity }

    private var selectedUseCase: UseCase? = null

    override fun onCreateView(
            inflater: LayoutInflater,
            container: ViewGroup?,
            savedInstanceState: Bundle?
    ) : View? {
        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_select_attach, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        // Get the use cases from the server
        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken: String?, _: String?, ex: AuthorizationException? ->
            if (accessToken == null) {
                Log.e(fTag, "[onViewCreated] got null access token, exception: $ex")
            } else {
                fetchUseCases(accessToken, view)
            }
        }

        // Set up the selection listener for the use case spinner
        attachment_list_view.setOnItemClickListener { adapterView, _, pos, _ ->

            val useCase = adapterView?.getItemAtPosition(pos) as? UseCase

            selectedUseCase = if (useCase != null) {
                Log.v(fTag, "[onItemSelected] User selected use case $useCase")
                (attachment_list_view.adapter as UseCaseAdapter).selectedUseCasePos = pos
                (attachment_list_view.adapter as UseCaseAdapter).notifyDataSetChanged()
                useCase
            } else {
                Log.e(fTag, "[onItemSelected] User indicated position $pos but adapter or"
                        + " item was null")
                parent.snackbar_layout.snackbar("Could not fetch your chosen use case")
                null
            }
        }

        register_button.setOnClickListener { _ -> registerRobotAndUseCase(selectedUseCase) }

        cancel_button.setOnClickListener { _ ->
            parent.startActivity(parent.intentFor<MainActivity>().clearTask())
        }
    }

    /**
     * Fetches all the use cases that the user can choose from.
     *
     * @param token the access token for the current session
     * @param view the parent view to inflate the use cases in
     */
    private fun fetchUseCases(token: String, view: View) {
        Log.v(fTag, "Getting use cases")
        parent.restApiService
                .getAllUseCases(token)
                .enqueue(object : Callback<List<UseCase>> {

            override fun onResponse(call: Call<List<UseCase>>, response: Response<List<UseCase>>) {
                list_view_progress_bar.visibility = View.GONE

                // Extract the use cases from the response, if it was successful
                val useCases: List<UseCase>? = response.body()
                if (response.isSuccessful && useCases != null) {
                    Log.v(fTag, "[fetchUseCases] Successfully got use cases: $useCases")

                    // Set up the listView with the downloaded use cases
                    attachment_list_view.adapter = UseCaseAdapter(view.context, useCases)
                    attachment_list_view.visibility = View.VISIBLE

                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.e(fTag, "[fetchUseCases] response was unsuccessful or body was null;"
                            + " error: $error")
                    if (error != null) {
                        parent.snackbar_layout.snackbar(error)
                    }
                }
            }

            override fun onFailure(call: Call<List<UseCase>>, t: Throwable) {
                val errorMsg = t.message
                Log.e(fTag, "[fetchUseCases] FAILED, got error: $errorMsg")
                if (errorMsg != null) {
                    parent.snackbar_layout.snackbar(errorMsg)
                }
            }

        })

    }

    /**
     * Registers the robot with ID [RegisterRobotActivity.robotId] with the chosen use case
     * to the user.
     *
     * @param useCase the use case chosen by the user for the robot, or null
     */
    private fun registerRobotAndUseCase(useCase: UseCase?) {

        if (useCase == null) {
            Log.e(fTag, "registerUseCase got null use case")
            return
        }

        use_case_registration_progress_bar.visibility = View.VISIBLE
        register_button.isEnabled = false
        cancel_button.isEnabled = false

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken: String?, _: String?, ex: AuthorizationException? ->
            // TODO am I supposed to use the accessToken or idToken (aka _)
            if (accessToken == null) {
                Log.e(fTag, "[registerRobotAndUseCase] got null access token, exception: $ex")
            } else {
                parent.restApiService
                        .registerRobot(
                                parent.robotId,
                                accessToken,
                                RobotRegistrationBody(parent.robotNickname, useCase.name))
                        .enqueue(object : Callback<ResponseBody> {

                            override fun onResponse(
                                    call: Call<ResponseBody>,
                                    response: Response<ResponseBody>) {

                                use_case_registration_progress_bar.visibility = View.GONE
                                register_button.isEnabled = true
                                cancel_button.isEnabled = true

                                if (response.isSuccessful) {
                                    context?.toast("Successfully registered the robot")
                                    parent.startFragment(ConfigureRobotFragment())
                                } else {
                                    val error = RestApiService.extractErrorFromResponse(response)
                                    Log.d(fTag, "[registerRobotAndUseCase] Got unsuccessful "
                                            + "response when registering robot and use case: $error")
                                    if (error != null) {
                                        parent.snackbar_layout.snackbar(error)
                                    }
                                }

                            }

                            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                                use_case_registration_progress_bar.visibility = View.GONE
                                register_button.isEnabled = true
                                cancel_button.isEnabled = true
                                val error = t.message
                                Log.e(fTag, "[registerRobotAndUseCase] FAILED, got error: $error")
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }
                            }
                        })
            }
        }
    }
}
