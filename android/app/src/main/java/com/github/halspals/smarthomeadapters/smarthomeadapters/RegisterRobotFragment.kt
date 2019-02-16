package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.text.Editable
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.inputmethod.InputMethodManager
import kotlinx.android.synthetic.main.fragment_register_robot.*
import okhttp3.ResponseBody
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * The second screen in the registration wizard, allowing the user to enter the robot nickname and then sending
 * the registration request.
 */
class RegisterRobotFragment : Fragment() {

    private val fTag = "RegisterRobotFragment"

    private lateinit var parent: RegisterRobotActivity

    override fun onAttach(context: Context?) {
        super.onAttach(context)
        parent = context as RegisterRobotActivity
    }

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_register_robot, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        register_button.setOnClickListener { _ ->

            // Dismiss the keyboard
            val inputMethodManager = activity?.getSystemService(Context.INPUT_METHOD_SERVICE) as? InputMethodManager
            val currentView = activity?.currentFocus ?: view
            inputMethodManager?.hideSoftInputFromWindow(currentView.windowToken, 0)

            registration_progress_bar.visibility = View.VISIBLE
            register_button.isEnabled = false

            // Make the registration call
            registerRobot(nickname_input.text)
        }

        id_text_view.text = getString(R.string.scanned_id_textiew, parent.robotId)

    }

    /**
     * Reads the nickname which has been input and if appropriate makes a call to the REST API to register the robot.
     *
     * @param nickname the nickname which the user has entered
     */
    private fun registerRobot(nickname: Editable) {
        if (nickname.isEmpty()) {
            nickname_input.error = "Enter a nickname for the robot"
        } else {
            parent
                .restApiService
                .registerRobot(parent.robotId, nickname.toString())
                .enqueue(object : Callback<ResponseBody> {

                    override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {

                        // Indicate to the user that the api call has been finished
                        registration_progress_bar.visibility = View.GONE
                        register_button.isEnabled = true

                        if (response.isSuccessful) {
                            // The robot was registered; move on to the next step of the wizard
                            context?.toast("Successfully registered the robot")
                            parent.startFragment(SelectUseCaseFragment())
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "registerRobot got unsuccessful response, error: $error")
                            // TODO display error to user
                        }
                    }

                    override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                        registration_progress_bar.visibility = View.GONE
                        register_button.isEnabled = true

                        val errorMsg = t.message
                        Log.e(fTag, "registerRobot FAILED, got error: $errorMsg")
                        // TODO display error to user
                    }
                })
        }
    }
}
