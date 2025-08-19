import React, { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Menu,
  MenuItem,
  Typography,
  Stack,
  TablePagination,
} from '@mui/material';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import AlertStatusChip from '../common/AlertStatusChip';
import { AlertStatus, AlertType } from '../../types/enums';
import { formatAlertType, formatPrice, formatDateTime } from '../../utils/formatters';

interface Alert {
  id: string;
  stockSymbol: string;
  stockName: string;
  alertType: AlertType;
  thresholdPrice?: number;
  currentPrice?: number;
  status: AlertStatus;
  createdAt: string;
}

interface AlertsTableProps {
  alerts: Alert[];
  onEdit: (alert: Alert) => void;
  onDelete: (alertId: string) => void;
}

const AlertsTable: React.FC<AlertsTableProps> = ({ alerts, onEdit, onDelete }) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedAlert, setSelectedAlert] = useState<Alert | null>(null);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>, alert: Alert) => {
    setAnchorEl(event.currentTarget);
    setSelectedAlert(alert);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedAlert(null);
  };

  const handleEdit = () => {
    if (selectedAlert) {
      onEdit(selectedAlert);
    }
    handleMenuClose();
  };

  const handleDelete = () => {
    if (selectedAlert) {
      onDelete(selectedAlert.id);
    }
    handleMenuClose();
  };

  const handleChangePage = (_: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const paginatedAlerts = alerts.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);

  return (
    <Paper>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Stock</TableCell>
              <TableCell>Alert Type</TableCell>
              <TableCell>Threshold Price</TableCell>
              <TableCell>Current Price</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Created</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedAlerts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Typography color="text.secondary" sx={{ py: 4 }}>
                    {alerts.length === 0 ? 'No alerts created yet' : 'No alerts match the current filter'}
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              paginatedAlerts.map((alert) => (
                <TableRow key={alert.id} hover>
                  <TableCell>
                    <Stack>
                      <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
                        {alert.stockSymbol}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {alert.stockName}
                      </Typography>
                    </Stack>
                  </TableCell>
                  <TableCell>
                    {formatAlertType(alert.alertType)}
                  </TableCell>
                  <TableCell>
                    {alert.thresholdPrice ? formatPrice(alert.thresholdPrice) : '-'}
                  </TableCell>
                  <TableCell>
                    {alert.currentPrice ? formatPrice(alert.currentPrice) : '-'}
                  </TableCell>
                  <TableCell>
                    <AlertStatusChip status={alert.status} />
                  </TableCell>
                  <TableCell>
                    {formatDateTime(new Date(alert.createdAt))}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={(e) => handleMenuClick(e, alert)}
                    >
                      <MoreVertIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      <TablePagination
        rowsPerPageOptions={[5, 10, 25]}
        component="div"
        count={alerts.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
      />

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleEdit}>
          <EditIcon fontSize="small" sx={{ mr: 1 }} />
          Edit Alert
        </MenuItem>
        <MenuItem onClick={handleDelete} sx={{ color: 'error.main' }}>
          <DeleteIcon fontSize="small" sx={{ mr: 1 }} />
          Delete Alert
        </MenuItem>
      </Menu>
    </Paper>
  );
};

export default AlertsTable;