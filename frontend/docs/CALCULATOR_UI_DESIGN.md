# Calculator-Style UI Design for Cash Flow Dialog

## Overview

The `AddCashFlowDialog` component has been redesigned with a calculator-style amount input for better user experience when entering amounts.

## Key Features

### 1. **Calculator-Style Amount Input**

- Large, borderless input field (similar to calculator screen)
- Background: Light gray (`bg-slate-100`) in light mode, dark gray (`bg-slate-900`) in dark mode
- Text: Extra large (4xl), bold, right-aligned
- Placeholder: Shows "0" when empty
- No default value (starts empty)
- Fully functional number input with keyboard support

### 2. **Layout Changes**

- Amount input moved to the top of the form (above Type and Category)
- Prominent position emphasizes the primary action (entering amount)
- Other fields remain in their original positions

## Components

### `AddCashFlowDialog.tsx` Changes

**UI Changes:**

- Replaced standard number input with calculator-style input field
- Moved amount input to the top of the form (above Type and Category)
- Large, borderless input with custom styling
- Right-aligned text for better readability
- Focus ring for accessibility

**Input Styling:**

```css
className="w-full bg-slate-100 dark:bg-slate-900 rounded-lg p-6 min-h-20 text-4xl font-bold tabular-nums text-right border-0 focus:outline-none focus:ring-2 focus:ring-ring"
```

**Form Integration:**

- Maintains react-hook-form integration
- Proper validation with Zod schema
- Value handling: empty string when no value, number when entered
- Parses input to number on change

## User Experience

### Before

- Standard number input field with border
- Default value of 0
- Manual keyboard input only

### After

- Large, prominent amount input at the top
- No border (cleaner look)
- Starts empty with "0" placeholder
- Keyboard input supported
- Right-aligned for better readability
- Focus ring for accessibility

## Technical Details

### Styling

- Uses Tailwind CSS utility classes
- Responsive design
- Dark mode support
- Consistent with shadcn/ui design system

### Form Integration

- Maintains react-hook-form integration
- Proper validation with Zod schema
- Hidden input for form submission
- Synced state between display and form value

### Accessibility

- Large, touch-friendly buttons
- Clear visual hierarchy
- Proper ARIA labels (inherited from Button component)
- Keyboard navigation support

## Future Enhancements

Potential improvements for future iterations:

- Add decimal point button for cents
- Add quick amount buttons (e.g., 100, 500, 1000)
- Add calculation operators (+, -, ×, ÷)
- Add memory functions (M+, M-, MR, MC)
- Add haptic feedback on mobile devices
- Add sound effects (optional)
- Add animation on button press
- Support for different currencies

## Testing

To test the calculator-style UI:

1. Open the cash flow page
2. Click "新增記錄" (Add Record)
3. Enter an amount using keyboard or mobile number pad
4. Verify the input displays correctly with large text
5. Test placeholder behavior (shows "0" when empty)
6. Submit the form and verify the amount is saved correctly

## Notes

- The amount input is positioned at the top of the form for prominence
- The input field has no border for a cleaner, calculator-like appearance
- Text is right-aligned and uses tabular numbers for better readability
- The form value is parsed to a number on change
- Empty input is treated as undefined (not 0)
