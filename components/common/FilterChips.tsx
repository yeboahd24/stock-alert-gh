import React from 'react';
import { Chip, Stack } from '@mui/material';

interface FilterChipsProps {
  filters: { label: string; value: string; active: boolean }[];
  onFilterChange: (value: string) => void;
}

const FilterChips: React.FC<FilterChipsProps> = ({ filters, onFilterChange }) => {
  return (
    <Stack direction="row" spacing={1} sx={{ flexWrap: 'wrap', gap: 1 }}>
      {filters.map((filter) => (
        <Chip
          key={filter.value}
          label={filter.label}
          variant={filter.active ? 'filled' : 'outlined'}
          color={filter.active ? 'primary' : 'default'}
          onClick={() => onFilterChange(filter.value)}
          clickable
        />
      ))}
    </Stack>
  );
};

export default FilterChips;