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
import com.google.zxing.integration.android.IntentIntegrator
import kotlinx.android.synthetic.main.activity_register_robot.*
import kotlinx.android.synthetic.main.fragment_qr.*
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * A fragment which upon a button click starts a QR scanner from the parent activity.
 * This is the first screen in the registration wizard.
 */
class QRFragment : Fragment() {

    private val fTag = "QRFragment"

    private lateinit var parent: RegisterRobotActivity

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_qr, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        parent = activity as RegisterRobotActivity

        scan_qr_button.setOnClickListener { _ ->
            // Start a QR scan from the parent using ZXING
            IntentIntegrator(activity)
                .setDesiredBarcodeFormats(IntentIntegrator.QR_CODE)
                .setPrompt("Scan robot QR code")
                .initiateScan()

        }

        register_button.setOnClickListener { _ ->
            // Dismiss the keyboard
            val inputMethodManager = parent.getSystemService(Context.INPUT_METHOD_SERVICE) as? InputMethodManager
            val currentView = parent.currentFocus ?: view
            inputMethodManager?.hideSoftInputFromWindow(currentView.windowToken, 0)

            register_button.isEnabled = false

            // Make the registration call
            registerRobot(manual_code_input.text, nickname_input.text)
        }

    }

    /**
     * Reads the nickname which has been input and if appropriate makes a call to the REST API to register the robot.
     *
     * @param code the robot's registration code
     * @param nickname the nickname which the user has entered
     */
    private fun registerRobot(code: Editable, nickname: Editable) {

        when {
            code.isEmpty() -> {
                manual_code_input.error = "Enter your robot's registration code"
                register_button.isEnabled = true
            }

            nickname.isEmpty() -> {
                nickname_input.error = "Enter a nickname for the robot"
                register_button.isEnabled = true
            }

            else -> {
                registration_progress_bar.visibility = View.VISIBLE

                parent
                    .restApiService
                    .registerRobot(code.toString(), nickname.toString())
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
                                parent.startFragment(SelectUseCaseFragment()) // TODO remove test code
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }

                            }
                        }

                        override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                            registration_progress_bar.visibility = View.GONE
                            register_button.isEnabled = true

                            val error = t.message
                            Log.e(fTag, "registerRobot FAILED, got error: $error")
                            if (error != null) {
                                parent.snackbar_layout.snackbar(error)
                            }
                        }
                    })
            }
        }
    }

    internal fun setRobotIdText(text: String) {
        manual_code_input.setText(text)
    }
}
