package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.google.zxing.integration.android.IntentIntegrator
import kotlinx.android.synthetic.main.fragment_qr.*

/**
 * A fragment which upon a button click starts a QR scanner from the parent activity.
 * This is the first screen in the registration wizard.
 */
class QRFragment : Fragment() {

    private val fTag = "QRFragment"

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_qr, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        scan_qr_button.setOnClickListener { _ ->
            // Start a QR scan from the parent using ZXING
            IntentIntegrator(activity)
                .setDesiredBarcodeFormats(IntentIntegrator.QR_CODE)
                .setPrompt("Scan robot QR code")
                .initiateScan()
        }

        manual_submit_button.setOnClickListener { _ ->
            (activity as RegisterRobotActivity).robotId = manual_code_input.text.toString()
            (activity as RegisterRobotActivity).startFragment(RegisterRobotFragment())
        }
    }
}
