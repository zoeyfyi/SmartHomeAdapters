package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.google.zxing.integration.android.IntentIntegrator
import kotlinx.android.synthetic.main.fragment_qr.*
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.intentFor

/**
 * A fragment which upon a button click starts a QR scanner from the parent activity.
 * This is the first screen in the registration wizard.
 */
class QRFragment : Fragment() {

    private val fTag = "QRFragment"

    private val parent by lazy { activity as RegisterRobotActivity }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_qr, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)


        camera_button.setOnClickListener { _ ->
            // Start a QR scan from the parent using ZXING
            IntentIntegrator(activity)
                .setDesiredBarcodeFormats(IntentIntegrator.QR_CODE)
                .setPrompt("Scan robot QR code")
                .initiateScan()

        }

        manual_entry_button.setOnClickListener { _ ->
            parent.startFragment(ManualEntryFragment())
        }

        cancel_button.setOnClickListener { _ ->
            parent.startActivity(parent.intentFor<MainActivity>().clearTask())
        }

    }
}
