package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigParameter
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigResult

/**
 * Provides an adapter to list [ConfigParameter]s in a grid of cards.
 *
 * @property context the context invoking the adapter
 * @property parameters the configuration parameters to list
 */
class ParameterAdapter(
        private val context: Context,
        private val parameters: List<ConfigParameter>
):  BaseAdapter() {

    private val configSettings = HashMap<String, String>()

    private val tag = "ParameterAdapter"

    init {
        for (param in parameters) {
            when (param.type) {
                ConfigParameter.BOOL_TYPE -> {
                    configSettings[param.name] = (param.details.default != 0).toString()
                }

                ConfigParameter.INT_TYPE -> {
                    configSettings[param.name] = param.details.default.toString()
                }

                else -> {
                    TODO("No other types expected")
                }
            }
        }
    }

    override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View {

        if (convertView != null) {
            // We've already been given a view which is set up -- use it
            return convertView
        }

        // Otherwise inflate the appropriate card and set up its fields
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val parameter = parameters[position]

        val view = when (parameter.type) {
            ConfigParameter.BOOL_TYPE -> {
                // The parameter is a boolean one -- inflate the bool card and set the default state
                val cardView = inflater.inflate(R.layout.view_bool_config_card, parent, false)
                val switch = cardView.findViewById<Switch>(R.id.input_switch)
                switch.isChecked = parameter.details.default != 0
                switch.setOnCheckedChangeListener { _, b ->
                    configSettings[parameter.name] = b.toString()
                }
                cardView
            }

            ConfigParameter.INT_TYPE -> {
                // The parameter is an integer one -- inflate the int card and set up the seek bar
                val cardView = inflater.inflate(R.layout.view_int_config_card, parent, false)

                val max = parameter.details.max
                val min = parameter.details.min
                val default = parameter.details.default

                val seekBar = cardView.findViewById<SeekBar>(R.id.input_seekbar)
                seekBar.max = max - min // manually add min later as api <26 doesn't support seekBar.setMin()
                seekBar.progress = default - min

                val seekBarValueTextView = cardView.findViewById<TextView>(R.id.seekbar_value_text_view)
                seekBarValueTextView.text = "$default"

                // Set up the change listener for the seekbar
                seekBar.setOnSeekBarChangeListener(object : SeekBar.OnSeekBarChangeListener {
                    override fun onStopTrackingTouch(p0: SeekBar?) {
                        // Not interested
                    }

                    override fun onStartTrackingTouch(p0: SeekBar?) {
                        // Not interested
                    }

                    override fun onProgressChanged(seekBar: SeekBar?, progress: Int, byUser: Boolean) {
                        // add back the min we subtracted earlier as the seekbar starts at 0
                        val newValue = (progress + min).toString()
                        seekBarValueTextView.text = newValue
                        configSettings[parameter.name] = newValue

                    }
                })

                cardView
            }

            else -> {
                TODO("No other types expected")
            }
        }

        // Regardless of the type of card, set up the name and description fields
        view.findViewById<TextView>(R.id.config_name_text_view).text = parameter.name
        view.findViewById<TextView>(R.id.config_explanation_text_view).text = parameter.description

        return view
    }

    override fun getCount(): Int = parameters.size

    override fun getItem(position: Int): Any = parameters[position]

    override fun getItemId(p0: Int): Long = 0L

    // TODO make more elegant solution than just stripping the space
    internal fun getConfigValuesSnapshot(): List<ConfigResult> {
        val snapshot = ArrayList<ConfigResult>()
        for (paramName in configSettings.keys) {
            snapshot.add(
                    ConfigResult(
                            paramName.replace(" ", ""),
                            configSettings[paramName]!!
                    )
            )
        }

        return snapshot
    }

}