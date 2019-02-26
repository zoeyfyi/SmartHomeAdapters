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
import kotlinx.android.synthetic.main.fragment_qr.*

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
            registerRobot(registration_code_input.text, nickname_input.text)
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
                registration_code_input.error = "Enter your robot's registration code"
                register_button.isEnabled = true
            }

            nickname.isEmpty() -> {
                nickname_input.error = "Enter a nickname for the robot"
                register_button.isEnabled = true
            }

            else -> {
                parent.robotNickname = nickname.toString()
                parent.robotId = code.toString()
                parent.startFragment(SelectUseCaseFragment())
            }
        }
    }

    internal fun setRobotIdText(text: String) {
        registration_code_input.setText(text)
    }
}
