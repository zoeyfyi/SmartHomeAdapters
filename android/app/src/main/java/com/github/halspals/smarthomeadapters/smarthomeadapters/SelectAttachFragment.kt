package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ArrayAdapter
import android.widget.BaseAdapter
import android.widget.CheckedTextView
import android.widget.TextView
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

            attachment_list_view.setItemChecked(pos, true)

            val useCase = adapterView?.getItemAtPosition(pos) as? UseCase

            selectedUseCase = if (useCase != null) {
                Log.v(fTag, "[onItemSelected] User selected use case $useCase")
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
                    Log.v(fTag, "[getAllUseCases] Successfully got list of ${useCases.size} "
                            + "use cases")

                    attachment_list_view.adapter = object : ArrayAdapter<UseCase>(context!!, 0, useCases) {

                        override fun getView(position: Int, convertView: View?, parent: ViewGroup): View {
                            val inflater =
                                view.context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
                            val ret = inflater.inflate(R.layout.attachment_list_item, parent, false)
                            ret.findViewById<TextView>(R.id.list_item_header).text = useCases[position].name
                            ret.findViewById<TextView>(R.id.list_item_description).text = useCases[position].name // TODO change obvs
                            return ret
                        }
                    }

                    attachment_list_view.visibility = View.VISIBLE

                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.e(fTag, "[getAllUseCases] response was unsuccessful or body was null;"
                            + " error: $error")
                    if (error != null) {
                        parent.snackbar_layout.snackbar(error)
                    }
                }
            }

            override fun onFailure(call: Call<List<UseCase>>, t: Throwable) {
                val errorMsg = t.message
                Log.e(fTag, "[getAllUseCases] FAILED, got error: $errorMsg")
                if (errorMsg != null) {
                    parent.snackbar_layout.snackbar(errorMsg)
                }
            }

        })

    }

    /**
     * WIP: Register the use case chosen to the robot being registered and set up.
     *
     * @param useCase the use case chosen by the user for the robot
     */
    private fun registerRobotAndUseCase(useCase: UseCase?) {

        if (useCase == null) {
            Log.e(fTag, "registerUseCase got null use case")
            return
        }

        use_case_registration_progress_bar.visibility = View.VISIBLE
        register_button.isEnabled = false
        cancel_button.isEnabled = false
        parent.chosenUseCase = useCase

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
                                RobotRegistrationBody(parent.robotNickname, useCase.name.toLowerCase()))
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
