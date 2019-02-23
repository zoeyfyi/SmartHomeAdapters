package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.AdapterView
import android.widget.ArrayAdapter
import android.widget.BaseAdapter
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigDetails
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigParameter
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.UseCase
import kotlinx.android.synthetic.main.activity_register_robot.*
import kotlinx.android.synthetic.main.fragment_select_use_case.*
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * The (optional) third screen in the robot registration wizard, allowing the user to choose a use case for the robot.
 */
class SelectUseCaseFragment : Fragment() {


    private val fTag = "SelectUseCaseFragment"

    private lateinit var parent: RegisterRobotActivity

    private var selectedUseCase: UseCase? = null

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_select_use_case, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        parent = activity as RegisterRobotActivity

        // Get the use cases from the server
        Log.v(fTag, "Getting use cases")
        parent.restApiService.getAllUseCases().enqueue(object : Callback<List<UseCase>> {

            override fun onResponse(call: Call<List<UseCase>>, response: Response<List<UseCase>>) {
                spinner_progress_bar.visibility = View.GONE

                // Extract the use cases from the response, if it was successful
                // TODO the spinner should not be set up if the call fails once the /usecases endpoint is working
                val useCases: List<UseCase>? = response.body()
                if (response.isSuccessful && useCases != null) {
                    Log.v(fTag, "[getAllUseCases] Successfully got list of ${useCases.size} use cases")

                    spinner.adapter = object : BaseAdapter() {
                        override fun getCount(): Int {
                            return useCases.size
                        }

                        override fun getItemId(p0: Int): Long {
                            return 0L
                        }

                        override fun getItem(position: Int): Any {
                            return useCases[position]
                        }

                        override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View {
                            if (convertView != null) {
                                return convertView
                            }

                            val inflater =
                                view.context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
                            val ret = inflater.inflate(android.R.layout.simple_spinner_item, parent, false)
                            ret.findViewById<TextView>(android.R.id.text1).text = useCases[position].name
                            return ret
                        }
                    }

                    spinner.visibility = View.VISIBLE
                } else {

                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.e(fTag, "[getAllUseCases] response was unsuccessful or body was null; error: $error")
                    if (error != null) {
                        parent.snackbar_layout.snackbar(error)
                    }
                    /*
                    TODO REMOVE TEST CODE.
                    Currently the app uses hard-coded alternatives for configuration so that it can be tested
                    while the end points are still being developed.
                    TODO THIS SHOULD BE REMOVED!
                     */

                    val testUseCases = listOf(
                        UseCase(
                            "1",
                            "Thermostat",
                            listOf(
                                ConfigParameter(
                                    "10 Degrees Point", "Set the angle corresponding to the temperature" +
                                            "10 degrees Celsius.", "int", ConfigDetails(75, 0, 90)
                                ),
                                ConfigParameter(
                                    "30 Degrees Point", "Set the angle corresponding to the temperature" +
                                            "30 degrees Celsius.", "int", ConfigDetails(145, 90, 180)
                                ),
                                ConfigParameter(
                                    "Test bool param", "just for testing; default ON", "bool",
                                    ConfigDetails(1, 0, 1)
                                )
                            )
                        ),
                        UseCase(
                            "2",
                            "Switch",
                            listOf(
                                ConfigParameter(
                                    "Off Angle", "Set the angle for turning off the light",
                                    "int", ConfigDetails(80, 0, 90)
                                ),
                                ConfigParameter(
                                    "On Angle", "Set the angle for turning on the light",
                                    "int", ConfigDetails(165, 90, 180)
                                ),
                                ConfigParameter(
                                    "Test bool param", "just for testing; default OFF", "bool",
                                    ConfigDetails(0, 0, 1)
                                )
                            )
                        )
                    )

                    spinner.adapter = object : BaseAdapter() {
                        override fun getCount(): Int {
                            return testUseCases.size
                        }

                        override fun getItemId(p0: Int): Long {
                            return 0L
                        }

                        override fun getItem(position: Int): Any {
                            return testUseCases[position]
                        }

                        override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View {
                            val inflater =
                                view.context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
                            val ret = inflater.inflate(android.R.layout.simple_spinner_item, parent, false)
                            ret.findViewById<TextView>(android.R.id.text1).text = testUseCases[position].name
                            return ret
                        }
                    }

                    spinner.visibility = View.VISIBLE


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

        // Set up the selection listener for the use case spinner
        spinner.onItemSelectedListener = object: AdapterView.OnItemSelectedListener {

            override fun onNothingSelected(p0: AdapterView<*>?) {
                selectedUseCase = null
            }

            override fun onItemSelected(adapter: AdapterView<*>?, view: View?, pos: Int, p3: Long) {

                val useCase = adapter?.getItemAtPosition(pos) as? UseCase

                selectedUseCase = if (useCase != null) {
                    Log.v(fTag, "[onItemSelected] User selected use case $useCase")
                    useCase
                } else {
                    Log.e(fTag, "[onItemSelected] User indicated position $pos but adapter or item was null")
                    parent.snackbar_layout.snackbar("Could not fetch your chosen use case")
                    null
                }
            }

        }

        // Set up the click listener for if the user wants to finish the registration wizard without use case selection
        // or configuration
        finish_early_button.setOnClickListener { _ -> activity?.finish() }

        set_usecase_button.setOnClickListener { _ -> registerUseCase(selectedUseCase) }

    }

    /**
     * WIP: Register the use case chosen to the robot being registered and set up.
     * TODO: Remove temporary error-ignoring code.
     *
     * @param useCase the use case chosen by the user for the robot
     */
    private fun registerUseCase(useCase: UseCase?) {

        if (useCase == null) {
            Log.e(fTag, "registerUseCase got null use case")
            return
        }

        use_case_registration_progress_bar.visibility = View.VISIBLE
        parent.chosenUseCase = useCase

        parent.restApiService.registerUseCaseToRobot(parent.robotId, useCase.id).enqueue(object : Callback<ResponseBody> {
            override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {

                use_case_registration_progress_bar.visibility = View.GONE

                if (response.isSuccessful) {
                    context?.toast("Successfully set up the use case")
                    parent.startFragment(ConfigureRobotFragment())
                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.d(fTag, "[registerUseCase] Got unsuccessful response when registering use case: $error")
                    if (error != null) {
                        parent.snackbar_layout.snackbar(error)
                    }
                    parent.startFragment(ConfigureRobotFragment()) // TODO THIS IS TEMP ONLY WHILE THE ENDPOINT IS NOT UP
                }

            }

            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                use_case_registration_progress_bar.visibility = View.GONE
                val error = t.message
                Log.e(fTag, "[getAllUseCases] FAILED, got error: $error")
                if (error != null) {
                    parent.snackbar_layout.snackbar(error)
                }
            }
        })
    }


}
