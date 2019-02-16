package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.AdapterView
import android.widget.ArrayAdapter
import kotlinx.android.synthetic.main.fragment_select_use_case.*
import okhttp3.ResponseBody
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * The (optional) third screen in the robot registration wizard, allowing the user to choose a use case for the robot.
 */
class SelectUseCaseFragment : Fragment() {


    private val fTag = "SelectUseCaseFragment"

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_select_use_case, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        // Get the use cases from the server
        Log.v(fTag, "Getting use cases")
        (activity as RegisterRobotActivity).restApiService.getAllUseCases().enqueue(object : Callback<List<String>> {

            override fun onResponse(call: Call<List<String>>, response: Response<List<String>>) {

                spinner_progress_bar.visibility = View.GONE

                // Extract the use cases from the response, if it was successful
                // TODO the spinner should not be set up if the call fails once the /usecases endpoint is working
                // TODO present user with errors throughout, allow them to try again/refresh
                val body = response.body()
                val spinnerContents: List<String> = when {
                    !response.isSuccessful -> {
                        val error = RestApiService.extractErrorFromResponse(response)
                        spinner_text_view.visibility = View.INVISIBLE
                        Log.d(fTag, "[getAllUseCases] Got unsuccessful response when loading use cases: $error")
                        listOf("Switch", "Thermostat", "Botlock") // TODO remove these when /usecases works
                    }

                    body == null -> {
                        spinner_text_view.visibility = View.INVISIBLE
                        Log.w(fTag, "[getAllUseCases] Response was successful but body is null")
                        listOf("Switch", "Thermostat", "Botlock") // TODO remove these when /usecases works
                    }

                    else -> {
                        // response was succesful and body is non-null
                        context?.toast("Successfully retrieved list of use cases")
                        body
                    }
                }

                // Set the spinner to be visible and add an adapter with the use cases retrieved
                spinner.adapter = ArrayAdapter<String>(
                    view.context,
                    android.R.layout.simple_spinner_item,
                    spinnerContents)
                spinner.visibility = View.VISIBLE
            }

            override fun onFailure(call: Call<List<String>>, t: Throwable) {
                val errorMsg = t.message
                Log.e(fTag, "[getAllUseCases] FAILED, got error: $errorMsg")
                // TODO present the error to the user
            }
        })

        // Set up the selection listener for the use case spinner
        spinner.onItemSelectedListener = object: AdapterView.OnItemSelectedListener {

            override fun onNothingSelected(p0: AdapterView<*>?) {
                TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
            }

            override fun onItemSelected(adapter: AdapterView<*>?, view: View?, pos: Int, p3: Long) {

                spinner.isEnabled = false

                // The user has chosen a use case -- fetch it and and invoke [registerUseCase].

                val useCase: String? = adapter?.getItemAtPosition(pos) as? String

                if (useCase != null) {
                    Log.v(fTag, "[onItemSelected] User chose use case $useCase")
                    registerUseCase(useCase)
                } else {
                    Log.e(fTag, "[onItemSelected] User indicated position $pos but adapter or item was null")
                }
            }

        }

        // Set up the click listener for if the user wants to finish the registration wizard without use case selection
        // or configuration
        finish_early_button.setOnClickListener { _ -> activity?.finish() }

    }

    /**
     * WIP: Register the use case chosen to the robot being registered and set up.
     * TODO: Displaying errors to users, removing temporary error-ignoring code.
     *
     * @param useCase the use case chosen by the user for the robot
     */
    private fun registerUseCase(useCase: String) {
        val parent = activity as RegisterRobotActivity
        use_case_registration_progress_bar.visibility = View.VISIBLE

        parent.restApiService.registerUseCaseToRobot(parent.robotId, useCase).enqueue(object : Callback<ResponseBody> {
            override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {

                use_case_registration_progress_bar.visibility = View.GONE

                if (response.isSuccessful) {
                    context?.toast("Successfully set up the use case")
                    parent.startFragment(ConfigureRobotFragment())
                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.d(fTag, "[registerUseCase] Got unsuccessful response when registering use case: $error")
                    parent.startFragment(ConfigureRobotFragment()) // TODO THIS IS TEMP ONLY WHILE THE ENDPOINT IS NOT UP
                }

            }

            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                use_case_registration_progress_bar.visibility = View.GONE
                val errorMsg = t.message
                Log.e(fTag, "[getAllUseCases] FAILED, got error: $errorMsg")
            }
        })
    }


}
